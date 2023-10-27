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
	"math"
)

//=============================================================================

func GetFilteringAnalysis(tx *gorm.DB, tsId uint, params *FilteringParams) (*FilteringResponse, error) {

	//--- Get trading system

	ts, err := db.GetTradingSystemById(tx, tsId)
	if err != nil {
		return nil, err
	}

	if ts == nil {
		return nil, req.NewRequestError("Missing trading system with id=%v", tsId)
	}

	//--- Get filtering config (if not provided)

	config, err := getFilteringConfig(tx, tsId, params)
	if err != nil {
		return nil, err
	}

	//--- Get instrument for costs

	inst, err := db.GetInstrumentById(tx, ts.InstrumentId)
	if err != nil {
		return nil, err
	}

	res := &FilteringResponse{}
	res.TradingSystem.Id   = ts.Id
	res.TradingSystem.Name = ts.Name
	res.FilteringConfig    = *config

	err =runFiltering(tx, float64(inst.Cost), res)

	return res, err
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func getFilteringConfig(tx *gorm.DB, tsId uint, params *FilteringParams) (*FilteringConfig, error) {
	if params.NoConfig {
		config := &FilteringConfig{}

		tsf, err := db.GetTsFilteringById(tx, tsId)
		if err != nil {
			return nil, err
		}

		//--- A trading system may not have a config attached

		if tsf != nil {
			config = recordToParams(tsf)
		}

		return config, nil
	} else {
		return &params.FilteringConfig, nil
	}
}

//=============================================================================

func recordToParams(r *db.TsFiltering) *FilteringConfig {
	fp := &FilteringConfig{}
	fp.LongShort.Enabled      = r.LsEnabled
	fp.LongShort.LongPeriod   = r.LsLongPeriod
	fp.LongShort.ShortPeriod  = r.LsShortPeriod
	fp.LongShort.Threshold    = r.LsThreshold
	fp.LongShort.ShortPosPerc = r.LsShortPosPerc
	fp.EquityAverage.Enabled  = r.MaEnabled
	fp.EquityAverage.Days     = r.MaDays

	return fp
}

//=============================================================================

func paramsToRecord(fp *FilteringParams) *db.TsFiltering {
	return nil
}

//=============================================================================

func runFiltering(tx *gorm.DB, cost float64, res *FilteringResponse) error {

	filters, err := createFilterChain(&res.FilteringConfig)
	if err != nil {
		return err
	}

	//--- Retrieve profits from database

	diList, err := db.FindDailyInfoByTsId(tx, res.Id)
	if err != nil {
		return err
	}

	//--- Creates slices

	size := len(diList)

	e := &res.Equities
	e.Days              = make([]int, size)
	e.UnfilteredProfit  = make([]float64, size)
	e.FilteredProfit    = make([]float64, size)
	e.UnfilteredDrawdown= make([]float64, size)
	e.FilteredDrawdown  = make([]float64, size)
	e.Average           = make([]float64, size)

	currUnfProfit := 0.0
	currFilProfit := 0.0

	maSum  := 0.0
	maDays := res.EquityAverage.Days

	for i, di := range diList {
		currCost      := di.OpenProfit - cost * math.Abs(float64(di.NumTrades * di.Position))
		currUnfProfit += currCost
		maSum         += currUnfProfit

		e.Days[i]             = di.Day
		e.UnfilteredProfit[i] = currUnfProfit

		if i < maDays -1 || maDays==0 {
			e.Average[i] = 0
		} else {
			if i>maDays -1 {
				maSum -= e.UnfilteredProfit[i-maDays]
			}
			e.Average[i] = maSum/float64(maDays)
		}

		if applyFilters(filters, e, i) {
			currFilProfit += currCost
		}

		e.FilteredProfit[i]   = currFilProfit
	}

	core.CalcDrawDown(&e.UnfilteredProfit, &e.UnfilteredDrawdown)
	core.CalcDrawDown(&e.FilteredProfit,   &e.FilteredDrawdown)

	return nil
}

//=============================================================================

func applyFilters(list *[]EquityFilter, eq *Equities, index int) bool {

	for _, ef := range *list {
		if !ef.compute(eq, index) {
			return false
		}
	}

	return true
}

//=============================================================================
