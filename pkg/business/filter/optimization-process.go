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
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"log/slog"
	"time"
)

//=============================================================================

const MaxResultSize = 1000

//=============================================================================
//===
//=== OptimizationProcess
//===
//=============================================================================

type OptimizationProcess struct {
	ts            *db.TradingSystem
	data          *[]db.DailyInfo
	params        *OptimizationRequest
	info          *OptimizationInfo
	runComparator *RunComparator
	stopping      bool
}

//=============================================================================

func (f *OptimizationProcess) Start() error {
	err := f.params.Validate()
	if err != nil {
		return err
	}

	//--- Start optimization process

	field := f.params.FieldToOptimize

	f.runComparator = NewRunComparator(field)

	f.info = NewOptimizationInfo(MaxResultSize, f.runComparator.compare)

	f.info.MaxSteps          = f.params.StepsCount()
	f.info.FieldToOptimize   = field
	f.info.Filters.PosProfit = f.params.EnablePosProfit
	f.info.Filters.OldVsNew  = f.params.EnableOldNew
	f.info.Filters.WinPerc   = f.params.EnableWinPerc
	f.info.Filters.EquVsAvg  = f.params.EnableEquAvg

	go f.generate()

	return nil
}

//=============================================================================

func (f *OptimizationProcess) Stop() {
	slog.Info("Stop: Stopping optimization process", "tsId", f.ts.Id)
	f.stopping = true
}

//=============================================================================

func (f *OptimizationProcess) GetInfo() *OptimizationInfo {
	return f.info
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func (f *OptimizationProcess) generate() {
	slog.Info("generate: Started", "tsId", f.ts.Id, "tsName", f.ts.Name, "combine", f.params.CombineFilters)

	if f.params.CombineFilters {
		f.generateCombined()
	} else {
		f.generateNotCombined()
	}

	f.info.EndTime = time.Now()
	f.info.Status  = OptimStatusComplete
	slog.Info("Complete.")
}

//=============================================================================

func (f *OptimizationProcess) generateCombined() {
}

//=============================================================================

func (f *OptimizationProcess) generateNotCombined() {

	f.initValues()

	//--- Positive profit over a period

	if f.params.EnablePosProfit {
		slog.Info("generateNotCombined: Optimizing positive profit", "tsId", f.ts.Id, "tsName", f.ts.Name)

		filters := &TradingFilters{
			PosProEnabled: true,
		}

		for _, posProDays := range f.params.PosProDays.Steps() {
			filters.PosProDays = posProDays
			res := RunAnalysis(f.ts, filters, f.data)
			if f.addResult(TypePosProfit, posProDays, -1, -1, res) {
				return
			}
		}
	}

	//--- Old vs new periods

	if f.params.EnableOldNew {
		slog.Info("generateNotCombined: Optimizing old vs new periods", "tsId", f.ts.Id, "tsName", f.ts.Name)

		filters := &TradingFilters{
			OldNewEnabled: true,
		}

		for _, oldNewOldDays := range f.params.OldNewOldDays.Steps() {
			for _, oldNewNewDays := range f.params.OldNewNewDays.Steps() {
				for _, oldNewOldPerc := range f.params.OldNewOldPerc.Steps() {
					filters.OldNewOldDays = oldNewOldDays
					filters.OldNewNewDays = oldNewNewDays
					filters.OldNewOldPerc = oldNewOldPerc
					res := RunAnalysis(f.ts, filters, f.data)
					if f.addResult(TypeOldVsNew, oldNewOldDays, oldNewNewDays, oldNewOldPerc, res) {
						return
					}
				}
			}
		}
	}

	//--- Winning % greater than

	if f.params.EnableWinPerc {
		slog.Info("generateNotCombined: Optimizing winning percentage", "tsId", f.ts.Id, "tsName", f.ts.Name)

		filters := &TradingFilters{
			WinPerEnabled: true,
		}

		for _, winPerDays := range f.params.WinPercDays.Steps() {
			for _, winPerPerc := range f.params.WinPercPerc.Steps() {
				filters.WinPerDays = winPerDays
				filters.WinPerValue= winPerPerc
				res := RunAnalysis(f.ts, filters, f.data)
				if f.addResult(TypeWinPerc, winPerDays, -1, winPerPerc, res) {
					return
				}
			}
		}
	}

	//--- Equity vs its average

	if f.params.EnableEquAvg {
		slog.Info("generateNotCombined: Optimizing equity vs its average", "tsId", f.ts.Id, "tsName", f.ts.Name)

		filters := &TradingFilters{
			EquAvgEnabled: true,
		}

		for _, equAvgDays := range f.params.EquAvgDays.Steps() {
			filters.EquAvgDays = equAvgDays
			res := RunAnalysis(f.ts, filters, f.data)
			if f.addResult(TypeEquVsAvg, equAvgDays, -1, -1, res) {
				return
			}
		}
	}
}

//=============================================================================

func (f *OptimizationProcess) addResult(filter string, days int, newDays int, percentage int, far *AnalysisResponse) bool {
	f.info.CurrStep++

	r := &Run{
		FilterType: filter,
		Days      : days,
		NewDays   : newDays,
		Percentage: percentage,
	}

	r.NetProfit   = far.Summary.FilProfit
	r.AvgTrade    = far.Summary.FilAverageTrade
	r.MaxDrawdown = far.Summary.FilMaxDrawdown

	f.info.Results.Add(r)
	f.updateValues(&far.Summary)

	//--- Check if we have to stop the process

	if f.stopping {
		slog.Info("addResult: Got stop request")
	}

	return f.stopping
}

//=============================================================================

func (f *OptimizationProcess) initValues() {

	//--- Run without filters to get the baseline

	res := RunAnalysis(f.ts, &TradingFilters{}, f.data)

	f.info.BaseValue = f.runComparator.getValue(&res.Summary)
	f.info.BestValue = f.info.BaseValue
}

//=============================================================================

func (f *OptimizationProcess) updateValues(s *Summary) {
	value := f.runComparator.getValue(s)

	if f.info.BestValue < value {
		f.info.BestValue = value
	}
}

//=============================================================================
