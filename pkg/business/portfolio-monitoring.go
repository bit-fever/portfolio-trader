//=============================================================================
/*
Copyright © 2023 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package business

import (
	"github.com/bit-fever/core/req"
	"github.com/bit-fever/portfolio-trader/pkg/business/inout"
	"github.com/bit-fever/portfolio-trader/pkg/core"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
	"math"
	"sort"
	"time"
)

//=============================================================================

func GetPortfolioMonitoring(tx *gorm.DB, params *inout.PortfolioMonitoringParams) (*inout.PortfolioMonitoringResponse, error) {

	//--- Get list of trading systems and check length

	tsMap, err := db.GetTradingSystemsBySourceIdAsMap(tx, params.TsIds)

	if err != nil {
		return nil, err
	}

	if len(tsMap) != len(params.TsIds) {
		return nil, req.NewNotFoundError("Missing some trading systems (input:%v. found:%v)", len(params.TsIds), len(tsMap))
	}

	//--- Get trading systems daily data

	fromDay := calcFromDay(params.Period)
	idsArray:= calcIdsArrayFromSourceIds(tsMap)
	diList, err := db.FindDailyInfoFromDay(tx, idsArray, fromDay)

	if err != nil {
		return nil, err
	}

	diMap := buildSortedMapOfInfo(diList)
	res   := buildMonitoringResult(diMap, tsMap)
	buildTotalInfo(res)

	return res, nil
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func calcFromDay(period int) int {
	now := time.Now()
	ago := now.Add(time.Duration(-period) * time.Hour * 24)
	y,m,d := ago.Date()

	return y*10000 + int(m)*100 + d
}

//=============================================================================

func calcIdsArrayFromSourceIds(tsMap map[uint]*db.TradingSystem) []uint {
	var list []uint

	for k,_ := range tsMap {
		list = append(list, k)
	}

	return list
}

//=============================================================================

func buildSortedMapOfInfo(list *[]db.DailyInfo) *map[uint][]*db.DailyInfo {
	tsMap := map[uint][]*db.DailyInfo{}

	for _, di := range *list {
		diAux := di
		tsList, ok := tsMap[di.TradingSystemId]

		if !ok {
			tsList = []*db.DailyInfo{}
		}

		tsMap[di.TradingSystemId] = append(tsList, &diAux)
	}

	for _, list := range tsMap {
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].Day < list[j].Day
		})
	}

	return &tsMap
}

//=============================================================================

func buildMonitoringResult(diMap *map[uint][]*db.DailyInfo, tsMap map[uint]*db.TradingSystem) *inout.PortfolioMonitoringResponse {
	res := &inout.PortfolioMonitoringResponse{}

	if len(*diMap) != 0 {
		res.TradingSystems = make([]*inout.TradingSystemMonitoring, len(*diMap))
	} else {
		return res
	}

	i := 0
	for key, list := range *diMap {
		ts := tsMap[key]
		res.TradingSystems[i] = buildTradingSystemMonitoring(ts, list)
		i++
	}

	return res
}

//=============================================================================

func buildTradingSystemMonitoring(ts *db.TradingSystem, list []*db.DailyInfo) *inout.TradingSystemMonitoring {
	tsa := inout.NewTradingSystemMonitoring(len(list))
	tsa.Id   = ts.Id
	tsa.Name = ts.Name

	currRawProfit := 0.0
	currNetProfit := 0.0
	currTrades    := 0

	//--- build data for a single trading system

	for i, di := range list {
		currRawProfit += di.OpenProfit
		currNetProfit += di.OpenProfit - float64(ts.CostPerTrade) * math.Abs(float64(di.NumTrades * di.Position))
		currTrades += di.NumTrades

		tsa.Days[i]      = di.Day
		tsa.RawProfit[i] = currRawProfit
		tsa.NetProfit[i] = currNetProfit
		tsa.NumTrades[i] = currTrades
	}

	core.CalcDrawDown(&tsa.RawProfit, &tsa.RawDrawdown)
	core.CalcDrawDown(&tsa.NetProfit, &tsa.NetDrawdown)

	return tsa
}

//=============================================================================

type TotalInfo struct {
	rawProfit float64
	netProfit float64
	totTrades int
}

//-----------------------------------------------------------------------------

func buildTotalInfo(pm *inout.PortfolioMonitoringResponse) {
	daySum := map[int]*TotalInfo{}

	//--- Collect all days with associated sums

	for _, tsm := range (*pm).TradingSystems {
		for i, day := range tsm.Days {
			ds, ok := daySum[day]

			if !ok {
				ds = &TotalInfo{}
				daySum[day] = ds
			}

			ds.rawProfit += tsm.RawProfit[i]
			ds.netProfit += tsm.NetProfit[i]
			ds.totTrades += tsm.NumTrades[i]
		}
	}

	//--- Convert map into list and sort it

	var res []int

	for k, _ := range daySum {
		res = append(res, k)
	}

	sort.Ints(res)
	pm.Days= res

	//--- Loop on all days and build total arrays

	pm.RawProfit   = make([]float64, len(res))
	pm.NetProfit   = make([]float64, len(res))
	pm.RawDrawdown = make([]float64, len(res))
	pm.NetDrawdown = make([]float64, len(res))
	pm.NumTrades   = make([]int,     len(res))

	for i, day := range res {
		ds := daySum[day]

		pm.RawProfit[i] = ds.rawProfit
		pm.NetProfit[i] = ds.netProfit
		pm.NumTrades[i] = ds.totTrades
	}

	core.CalcDrawDown(&pm.RawProfit, &pm.RawDrawdown)
	core.CalcDrawDown(&pm.NetProfit, &pm.NetDrawdown)
}

//=============================================================================
