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
	"github.com/bit-fever/core/auth"
	"github.com/bit-fever/core/req"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func GetTradingFilters(tx *gorm.DB, c *auth.Context, tsId uint) (*db.TradingFilters, error) {
	_, err := getTradingSystem(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	return db.GetTradingFiltersByTsId(tx, tsId)
}

//=============================================================================

func SetTradingFilters(tx *gorm.DB, c *auth.Context, tsId uint, f *TradingFilters) error {
	_, err := getTradingSystem(tx, c, tsId)
	if err != nil {
		return err
	}

	db.SetTradingFilters(tx, &db.TradingFilters{
		TradingSystemId : tsId,
		EquAvgEnabled   : f.EquAvgEnabled,
		EquAvgDays      : f.EquAvgDays,
		PosProEnabled   : f.PosProEnabled,
		PosProWeeks     : f.PosProWeeks,
		WinPerEnabled   : f.WinPerEnabled,
		WinPerWeeks     : f.WinPerWeeks,
		WinPerValue     : f.WinPerValue,
		ShoLonEnabled   : f.ShoLonEnabled,
		ShoLonShortWeeks: f.ShoLonShortWeeks,
		ShoLonLongWeeks : f.ShoLonLongWeeks,
		ShoLonLongPerc  : f.ShoLonLongPerc,
	})

	return nil
}

//=============================================================================

func RunFilterAnalysis(tx *gorm.DB, c *auth.Context, tsId uint, far *FilterAnalysisRequest) (*FilterAnalysisResponse, error){
	ts, err := getTradingSystem(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	//--- Get filtering config (if not provided)

	if far.Filters == nil {
		filters, err := db.GetTradingFiltersByTsId(tx, tsId)
		if err != nil {
			return nil, err
		}

		far.Filters = convert(filters)
	}

	res := &FilterAnalysisResponse{}
	res.TradingSystem.Id   = ts.Id
	res.TradingSystem.Name = ts.Name
	res.Filters            = *far.Filters

	err = runFiltering(tx, ts, res)

	return res, err
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func getTradingSystem(tx *gorm.DB, c *auth.Context, tsId uint) (*db.TradingSystem, error){
	ts, err := db.GetTradingSystemById(tx, tsId)
	if err != nil {
		return nil, err
	}

	if ts == nil {
		return nil, req.NewNotFoundError("Trading system was not found: %v", tsId)
	}

	if ts.Username != c.Session.Username {
		return nil, req.NewForbiddenError("Trading system not owned by user: %v", tsId)
	}

	return ts, nil
}

//=============================================================================

func convert(f *db.TradingFilters) *TradingFilters {
	return &TradingFilters{
		EquAvgEnabled   : f.EquAvgEnabled,
		EquAvgDays      : f.EquAvgDays,
		PosProEnabled   : f.PosProEnabled,
		PosProWeeks     : f.PosProWeeks,
		WinPerEnabled   : f.WinPerEnabled,
		WinPerWeeks     : f.WinPerWeeks,
		WinPerValue     : f.WinPerValue,
		ShoLonEnabled   : f.ShoLonEnabled,
		ShoLonShortWeeks: f.ShoLonShortWeeks,
		ShoLonLongWeeks : f.ShoLonLongWeeks,
		ShoLonLongPerc  : f.ShoLonLongPerc,
	}
}

//=============================================================================

func runFiltering(tx *gorm.DB, ts *db.TradingSystem, res *FilterAnalysisResponse) error {

	//--- Retrieve profits from database

	diList, err := db.FindDailyInfoByTsId(tx, ts.Id)
	if err != nil {
		return err
	}

	//--- Creates slices

	_ = len(*diList)

	//e := &res.Equities
	//e.Days              = make([]int, size)
	//e.UnfilteredProfit  = make([]float64, size)
	//e.FilteredProfit    = make([]float64, size)
	//e.UnfilteredDrawdown= make([]float64, size)
	//e.FilteredDrawdown  = make([]float64, size)
	//e.Average           = make([]float64, size)
	//
	//currUnfProfit := 0.0
	//currFilProfit := 0.0
	//
	//maSum  := 0.0
	//maDays := res.EquityAverage.Days
	//
	//for i, di := range *diList {
	//	currCost      := di.OpenProfit - cost * math.Abs(float64(di.NumTrades * di.Position))
	//	currUnfProfit += currCost
	//	maSum         += currUnfProfit
	//
	//	e.Days[i]             = di.Day
	//	e.UnfilteredProfit[i] = currUnfProfit
	//
	//	if i < maDays -1 || maDays==0 {
	//		e.Average[i] = 0
	//	} else {
	//		if i>maDays -1 {
	//			maSum -= e.UnfilteredProfit[i-maDays]
	//		}
	//		e.Average[i] = maSum/float64(maDays)
	//	}
	//
	//	if applyFilters(filters, e, i) {
	//		currFilProfit += currCost
	//	}
	//
	//	e.FilteredProfit[i]   = currFilProfit
	//}
	//
	//core.CalcDrawDown(&e.UnfilteredProfit, &e.UnfilteredDrawdown)
	//core.CalcDrawDown(&e.FilteredProfit,   &e.FilteredDrawdown)

	return nil
}

//=============================================================================

//func applyFilters(list *[]EquityFilter, eq *Equities, index int) bool {
//
//	for _, ef := range *list {
//		if !ef.compute(eq, index) {
//			return false
//		}
//	}
//
//	return true
//}

//=============================================================================
