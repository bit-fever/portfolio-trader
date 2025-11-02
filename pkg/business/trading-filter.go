//=============================================================================
/*
Copyright © 2024 Andrea Carboni andrea.carboni71@gmail.com

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
	"github.com/tradalia/core/auth"
	"github.com/tradalia/portfolio-trader/pkg/business/filter"
	"github.com/tradalia/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func GetTradingFilters(tx *gorm.DB, c *auth.Context, tsId uint) (*db.TradingFilter, error) {
	_, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	return db.GetTradingFilterByTsId(tx, tsId)
}

//=============================================================================

func SetTradingFilters(tx *gorm.DB, c *auth.Context, tsId uint, f *filter.TradingFilter) error {
	_, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return err
	}

	tf := convert(f)
	tf.TradingSystemId = tsId

	return db.SetTradingFilter(tx, tf)
}

//=============================================================================

func RunFilterAnalysis(tx *gorm.DB, c *auth.Context, tsId uint, far *filter.AnalysisRequest) (*filter.AnalysisResponse, error){
	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	//--- Get filtering config (if not provided)
	var filters *db.TradingFilter

	if far.Filter == nil {
		filters, err = db.GetTradingFilterByTsId(tx, tsId)
		if err != nil {
			return nil, err
		}

	} else {
		filters = convert(far.Filter)
	}

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, far.StartDate, nil)
	if err != nil {
		return nil,err
	}

	res := filter.RunAnalysis(ts, filters, trades)

	return res, err
}

//=============================================================================

func StartFilterOptimization(tx *gorm.DB, c *auth.Context, tsId uint, oreq *filter.OptimizationRequest) error {
	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return err
	}

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, oreq.StartDate, nil)
	if err != nil {
		return err
	}

	err = oreq.Validate()
	if err != nil {
		return err
	}

	c.Log.Info("StartFilterOptimization: Starting optimization", "tsId", ts.Id, "tsName", ts.Name)
	filter.StartOptimization(ts, trades, oreq)

	return nil
}

//=============================================================================

func StopFilterOptimization(c *auth.Context, tsId uint) error {
	c.Log.Info("StopFilterOptimization: Stopping optimization", "tsId", tsId)
	err := filter.StopOptimization(tsId)

	return err
}

//=============================================================================

func GetFilterOptimizationInfo(c *auth.Context, tsId uint) (*filter.OptimizationResponse, error) {
	info := filter.GetOptimizationInfo(tsId)
	return filter.NewOptimizationResponse(info), nil
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func convert(f *filter.TradingFilter) *db.TradingFilter {
	return &db.TradingFilter{
		EquAvgEnabled   : f.EquAvgEnabled,
		EquAvgLen       : f.EquAvgLen,
		PosProEnabled   : f.PosProEnabled,
		PosProLen       : f.PosProLen,
		WinPerEnabled   : f.WinPerEnabled,
		WinPerLen       : f.WinPerLen,
		WinPerValue     : f.WinPerValue,
		OldNewEnabled   : f.OldNewEnabled,
		OldNewOldLen    : f.OldNewOldLen,
		OldNewOldPerc   : f.OldNewOldPerc,
		OldNewNewLen    : f.OldNewNewLen,
		TrendlineEnabled: f.TrendlineEnabled,
		TrendlineLen    : f.TrendlineLen,
		TrendlineValue  : f.TrendlineValue,
		DrawdownEnabled : f.DrawdownEnabled,
		DrawdownMin     : f.DrawdownMin,
		DrawdownMax     : f.DrawdownMax,
	}
}

//=============================================================================
