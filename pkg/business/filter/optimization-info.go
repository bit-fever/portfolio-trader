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

package filter

import (
	"github.com/bit-fever/portfolio-trader/pkg/business/filter/algorithm/optimization"
	"github.com/bit-fever/portfolio-trader/pkg/core"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"sync"
	"time"
)

//=============================================================================
//===
//=== Run
//===
//=============================================================================

type Run struct {
	Filter       *db.TradingFilter `json:"filter"`

	FitnessValue float64 `json:"fitnessValue"`
	NetProfit    float64 `json:"netProfit"`
	AvgTrade     float64 `json:"avgTrade"`
	MaxDrawdown  float64 `json:"maxDrawdown"`
	random       int
}

//=============================================================================
//===
//=== OptimizationInfo
//===
//=============================================================================

const OptimStatusIdle     = "idle"
const OptimStatusRunning  = "running"
const OptimStatusComplete = "complete"

type OptimizationInfo struct {
	sync.RWMutex
	CurrStep  uint
	MaxSteps  uint
	StartTime time.Time
	EndTime   time.Time
	Status    string
	results   *core.SortedResults

	StartDate       *time.Time
	BaseValue       float64
	BestValue       float64
	FieldToOptimize string

	Filter struct {
		PosProfit bool
		OldVsNew  bool
		WinPerc   bool
		EquVsAvg  bool
		Trendline bool
		Drawdown  bool
	}
}

//=============================================================================

func NewOptimizationInfo(maxResultSize int, field string, fc *optimization.FilterConfig,
						 steps uint, baseValue float64, startDate *time.Time) *OptimizationInfo {
	oi := &OptimizationInfo{}
	oi.CurrStep        = 0
	oi.StartTime       = time.Now()
	oi.Status          = OptimStatusRunning
	oi.results         = core.NewSortedResults(maxResultSize, runComparator)
	oi.BaseValue       = baseValue
	oi.BestValue       = baseValue
	oi.MaxSteps        = steps
	oi.FieldToOptimize = field
	oi.StartDate       = startDate

	oi.Filter.PosProfit = fc.EnablePosProfit
	oi.Filter.OldVsNew  = fc.EnableOldNew
	oi.Filter.WinPerc   = fc.EnableWinPerc
	oi.Filter.EquVsAvg  = fc.EnableEquAvg
	oi.Filter.Trendline = fc.EnableTrendline
	oi.Filter.Drawdown  = fc.EnableDrawdown

	return oi
}

//=============================================================================
//===
//=== Public methods
//===
//=============================================================================

func (oi *OptimizationInfo) GetRuns() []any {
	oi.Lock()
	defer oi.Unlock()

	if oi.results != nil {
		return oi.results.ToList()
	}

	return nil
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func (oi *OptimizationInfo) addResult(r *Run) {
	oi.Lock()
	defer oi.Unlock()

	oi.CurrStep++
	oi.results.Add(r)

	fv := r.FitnessValue

	if oi.BestValue < fv {
		oi.BestValue = fv
	}
}

//=============================================================================

func (oi *OptimizationInfo) isStatusComplete() bool {
	oi.Lock()
	defer oi.Unlock()

	if oi.CurrStep != oi.MaxSteps {
		return false
	}

	oi.EndTime = time.Now()
	oi.Status  = OptimStatusComplete

	return true
}

//=============================================================================
