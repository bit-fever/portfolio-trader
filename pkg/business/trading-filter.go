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
	"github.com/bit-fever/core/auth"
	"github.com/bit-fever/core/req"
	"github.com/bit-fever/portfolio-trader/pkg/business/filter"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func GetTradingFilters(tx *gorm.DB, c *auth.Context, tsId uint) (*db.TradingFilter, error) {
	_, err := getTradingSystem(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	return db.GetTradingFilterByTsId(tx, tsId)
}

//=============================================================================

func SetTradingFilters(tx *gorm.DB, c *auth.Context, tsId uint, f *filter.TradingFilter) error {
	_, err := getTradingSystem(tx, c, tsId)
	if err != nil {
		return err
	}

	db.SetTradingFilter(tx, &db.TradingFilter{
		TradingSystemId: tsId,
		EquAvgEnabled  : f.EquAvgEnabled,
		EquAvgLen      : f.EquAvgLen,
		PosProEnabled  : f.PosProEnabled,
		PosProLen      : f.PosProLen,
		WinPerEnabled  : f.WinPerEnabled,
		WinPerLen      : f.WinPerLen,
		WinPerValue    : f.WinPerValue,
		OldNewEnabled  : f.OldNewEnabled,
		OldNewOldLen   : f.OldNewOldLen,
		OldNewOldPerc  : f.OldNewOldPerc,
		OldNewNewLen   : f.OldNewNewLen,
	})

	return nil
}

//=============================================================================

func RunFilterAnalysis(tx *gorm.DB, c *auth.Context, tsId uint, far *filter.AnalysisRequest) (*filter.AnalysisResponse, error){
	ts, err := getTradingSystem(tx, c, tsId)
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

	diList, err := db.FindTradesByTsId(tx, ts.Id)
	if err != nil {
		return nil,err
	}

	res := filter.RunAnalysis(ts, filters, diList)

	return res, err
}

//=============================================================================

func StartFilterOptimization(tx *gorm.DB, c *auth.Context, tsId uint, far *filter.OptimizationRequest) error {
	ts, err := getTradingSystem(tx, c, tsId)
	if err != nil {
		return err
	}

	tradeList, err := db.FindTradesByTsId(tx, ts.Id)
	if err != nil {
		return err
	}

	c.Log.Info("StartFilterOptimization: Starting optimization", "tsId", ts.Id, "tsName", ts.Name)
	err = filter.StartOptimization(ts, tradeList, far)

	return err
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

func convert(f *filter.TradingFilter) *db.TradingFilter {
	return &db.TradingFilter{
		EquAvgEnabled : f.EquAvgEnabled,
		EquAvgLen     : f.EquAvgLen,
		PosProEnabled : f.PosProEnabled,
		PosProLen     : f.PosProLen,
		WinPerEnabled : f.WinPerEnabled,
		WinPerLen     : f.WinPerLen,
		WinPerValue   : f.WinPerValue,
		OldNewEnabled : f.OldNewEnabled,
		OldNewOldLen  : f.OldNewOldLen,
		OldNewOldPerc : f.OldNewOldPerc,
		OldNewNewLen  : f.OldNewNewLen,
	}
}

//=============================================================================
