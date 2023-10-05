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
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"github.com/bit-fever/portfolio-trader/pkg/tool"
	"gorm.io/gorm"
	"math"
	"sort"
	"time"
)

//=============================================================================

func GetPortfolioMonitoring(tx *gorm.DB, params *PortfolioMonitoringParams) (*PortfolioMonitoringResponse, error) {

	//--- Get list of trading systems and check length

	tsMap, err := db.GetTradingSystemsByIdAsMap(tx, params.TsIds)

	if err != nil {
		return nil, err
	}

	if len(tsMap) != len(params.TsIds) {
		return nil, tool.NewRequestError("Missing some trading systems (input:%v. found:%v)", len(params.TsIds), len(tsMap))
	}

	//--- Get trading systems daily data

	fromDay := calcFromDay(params.Period)
	diList, err := db.FindDailyInfoFromDay(tx, params.TsIds, fromDay)

	if err != nil {
		return nil, err
	}

	inMap, err := db.GetInstrumentsAsMap(tx)

	if err != nil {
		return nil, err
	}

	diMap := buildSortedMapOfInfo(diList)
	res   := buildMonitoringResult(diMap, &inMap, tsMap)
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

func buildSortedMapOfInfo(list *[]db.TsDailyInfo) *map[uint][]*db.TsDailyInfo {
	tsMap := map[uint][]*db.TsDailyInfo{}

	for _, di := range *list {
		diAux := di
		tsList, ok := tsMap[di.TradingSystemId]

		if !ok {
			tsList = []*db.TsDailyInfo{}
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

func buildMonitoringResult(diMap *map[uint][]*db.TsDailyInfo, inMap *map[uint]*db.Instrument, tsMap map[int]*db.TradingSystem) *PortfolioMonitoringResponse {
	res := &PortfolioMonitoringResponse{}

	if len(*diMap) != 0 {
		res.TradingSystems = make([]*TradingSystemMonitoring, len(*diMap))
	} else {
		return res
	}

	i := 0
	for key, list := range *diMap {
		ts := tsMap[int(key)]
		res.TradingSystems[i] = buildTradingSystemMonitoring(ts, list, inMap)
		i++
	}

	return res
}

//=============================================================================

func buildTradingSystemMonitoring(ts *db.TradingSystem, list []*db.TsDailyInfo, inMap *map[uint]*db.Instrument) *TradingSystemMonitoring {
	tsa := NewTradingSystemAnalysis(len(list))
	tsa.Id   = ts.Id
	tsa.Name = ts.Name

	currRawProfit := 0.0
	currNetProfit := 0.0
	currTrades    := 0

	cost := float64((*inMap)[ts.InstrumentId].Cost)

	//--- build data for a single trading system

	for i, di := range list {
		currRawProfit += di.OpenProfit
		currNetProfit += di.OpenProfit - cost * math.Abs(float64(di.NumTrades * di.Position))
		currTrades += di.NumTrades

		tsa.Days[i]      = di.Day
		tsa.RawProfit[i] = currRawProfit
		tsa.NetProfit[i] = currNetProfit
		tsa.NumTrades[i] = currTrades
	}

	calcDrawDown(&tsa.RawProfit, &tsa.RawDrawdown)
	calcDrawDown(&tsa.NetProfit, &tsa.NetDrawdown)

	return tsa
}

//=============================================================================

func calcDrawDown(equity *[]float64, drawDown *[]float64) {
	maxProfit    := 0.0
	currDrawDown := 0.0

	for i, currProfit := range *equity {
		if currProfit >= maxProfit {
			maxProfit = currProfit
			currDrawDown = 0
		} else {
			currDrawDown = currProfit - maxProfit
		}

		(*drawDown)[i] = currDrawDown
	}
}

//=============================================================================

type TotalInfo struct {
	rawProfit float64
	netProfit float64
	totTrades int
}

//-----------------------------------------------------------------------------

func buildTotalInfo(pm *PortfolioMonitoringResponse) {
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

	calcDrawDown(&pm.RawProfit, &pm.RawDrawdown)
	calcDrawDown(&pm.NetProfit, &pm.NetDrawdown)
}

//=============================================================================
