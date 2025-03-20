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
	"github.com/bit-fever/portfolio-trader/pkg/core"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
	"sort"
	"time"
)

//=============================================================================

func GetPortfolioMonitoring(tx *gorm.DB, params *PortfolioMonitoringParams) (*PortfolioMonitoringResponse, error) {

	//--- Get list of trading systems and check length

	tsMap, err := db.GetTradingSystemsByIdsAsMap(tx, params.TsIds)

	if err != nil {
		return nil, err
	}

	if len(tsMap) != len(params.TsIds) {
		return nil, req.NewNotFoundError("Missing some trading systems (input:%v. found:%v)", len(params.TsIds), len(tsMap))
	}

	//--- Get trading systems daily data

	fromTime := calcFromTime(params.Period)
	idsArray := calcIdsArrayFromSourceIds(tsMap)
	trades, err := db.FindTradesFromTime(tx, idsArray, fromTime)

	if err != nil {
		return nil, err
	}

	trMap := buildSortedMapOfInfo(trades)
	res   := buildMonitoringResult(trMap, tsMap)
	buildTotalInfo(res)

	return res, nil
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func calcFromTime(period int) time.Time {
	now := time.Now()

	return now.Add(time.Duration(-period) * time.Hour * 24)
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

func buildSortedMapOfInfo(list *[]db.Trade) *map[uint][]*db.Trade {
	trMap := map[uint][]*db.Trade{}

	for _, trade := range *list {
		trList, ok := trMap[trade.TradingSystemId]

		if !ok {
			trList = []*db.Trade{}
		}

		trMap[trade.TradingSystemId] = append(trList, &trade)
	}

	for _, list := range trMap {
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].ExitDate.Before(*list[j].ExitDate)
		})
	}

	return &trMap
}

//=============================================================================

func buildMonitoringResult(trMap *map[uint][]*db.Trade, tsMap map[uint]*db.TradingSystem) *PortfolioMonitoringResponse {
	res := &PortfolioMonitoringResponse{}

	if len(*trMap) != 0 {
		res.TradingSystems = make([]*TradingSystemMonitoring, len(*trMap))
	} else {
		return res
	}

	i := 0
	for key, list := range *trMap {
		ts := tsMap[key]
		res.TradingSystems[i] = buildTradingSystemMonitoring(ts, list)
		i++
	}

	return res
}

//=============================================================================

func buildTradingSystemMonitoring(ts *db.TradingSystem, list []*db.Trade) *TradingSystemMonitoring {
	tsa := NewTradingSystemMonitoring(len(list))
	tsa.Id   = ts.Id
	tsa.Name = ts.Name

	currRawProfit := 0.0
	currNetProfit := 0.0

	//--- build data for a single trading system

	for i, tr := range list {
		currRawProfit += tr.GrossProfit
		currNetProfit += tr.GrossProfit - float64(ts.CostPerOperation) * 2

		tsa.Time[i]        = *tr.ExitDate
		tsa.GrossProfit[i] = currRawProfit
		tsa.NetProfit[i]   = currNetProfit
	}

	core.CalcDrawDown(&tsa.GrossProfit, &tsa.GrossDrawdown)
	core.CalcDrawDown(&tsa.NetProfit,   &tsa.NetDrawdown)

	return tsa
}

//=============================================================================

type TotalInfo struct {
	grossProfit float64
	netProfit   float64
}

//-----------------------------------------------------------------------------

func buildTotalInfo(pm *PortfolioMonitoringResponse) {
	timeSum := map[time.Time]*TotalInfo{}

	//--- Collect all days with associated sums

	for _, tsm := range (*pm).TradingSystems {
		for i, t := range tsm.Time {
			ds, ok := timeSum[t]

			if !ok {
				ds = &TotalInfo{}
				timeSum[t] = ds
			}

			ds.grossProfit += tsm.GrossProfit[i]
			ds.netProfit   += tsm.NetProfit[i]
		}
	}

	//--- Convert map into list and sort it

	var res []time.Time

	for k, _ := range timeSum {
		res = append(res, k)
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Before(res[j])
	})

	pm.Time = res

	//--- Loop on all days and build total arrays

	pm.GrossProfit   = make([]float64, len(res))
	pm.NetProfit     = make([]float64, len(res))
	pm.GrossDrawdown = make([]float64, len(res))
	pm.NetDrawdown   = make([]float64, len(res))

	for i, day := range res {
		ds := timeSum[day]

		pm.GrossProfit[i] = ds.grossProfit
		pm.NetProfit[i]   = ds.netProfit
	}

	core.CalcDrawDown(&pm.GrossProfit, &pm.GrossDrawdown)
	core.CalcDrawDown(&pm.NetProfit,   &pm.NetDrawdown)
}

//=============================================================================
