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
	data          *[]db.Trade
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
	f.info.Filters.Trendline = f.params.EnableTrendline

	f.initValues()

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

	for ; !f.info.isStatusComplete(); {
		time.Sleep(time.Second * 1)
	}

	slog.Info("Complete.")
}

//=============================================================================

func (f *OptimizationProcess) generateCombined() {
}

//=============================================================================

func (f *OptimizationProcess) generateNotCombined() {
	if !f.generatePosProfit(){
		if !f.generateOldVsNew() {
			if !f.generateWinPerc() {
				if !f.generateEquVsAvg() {
					f.generateTrendline()
				}
			}
		}
	}
}

//=============================================================================

func (f *OptimizationProcess) generatePosProfit() bool {
	if f.params.EnablePosProfit {
		slog.Info("generatePosProfit: Optimizing positive profit", "tsId", f.ts.Id, "tsName", f.ts.Name)

		for _, posProLen := range f.params.PosProLen.Steps() {
			filter := &db.TradingFilter{
				PosProEnabled: true,
				PosProLen: posProLen,
			}

			go func() {
				res := RunAnalysis(f.ts, filter, f.data)
				f.addResult(TypePosProfit, posProLen, -1, -1, res)
			}()

			//--- Check if we have to stop the process

			if f.stopping {
				slog.Info("generatePosProfit: Got stop request")
				return true
			}
		}
	}

	return false
}

//=============================================================================

func (f *OptimizationProcess) generateOldVsNew() bool {
	if f.params.EnableOldNew {
		slog.Info("generateOldVsNew: Optimizing old vs new periods", "tsId", f.ts.Id, "tsName", f.ts.Name)

		for _, oldNewOldLen := range f.params.OldNewOldLen.Steps() {
			for _, oldNewNewLen := range f.params.OldNewNewLen.Steps() {
				for _, oldNewOldPerc := range f.params.OldNewOldPerc.Steps() {
					filter := &db.TradingFilter{
						OldNewEnabled: true,
						OldNewOldLen : oldNewOldLen,
						OldNewNewLen : oldNewNewLen,
						OldNewOldPerc: oldNewOldPerc,
					}

					go func() {
						res := RunAnalysis(f.ts, filter, f.data)
						f.addResult(TypeOldVsNew, oldNewOldLen, oldNewNewLen, oldNewOldPerc, res)
					}()

					//--- Check if we have to stop the process

					if f.stopping {
						slog.Info("generateOldVsNew: Got stop request")
						return true
					}
				}
			}
		}
	}

	return false
}

//=============================================================================

func (f *OptimizationProcess) generateWinPerc() bool {
	if f.params.EnableWinPerc {
		slog.Info("generateWinPerc: Optimizing winning percentage", "tsId", f.ts.Id, "tsName", f.ts.Name)

		for _, winPerLen := range f.params.WinPercLen.Steps() {
			for _, winPerPerc := range f.params.WinPercPerc.Steps() {
				filter := &db.TradingFilter{
					WinPerEnabled: true,
					WinPerLen    : winPerLen,
					WinPerValue  : winPerPerc,
				}

				go func(){
					res := RunAnalysis(f.ts, filter, f.data)
					f.addResult(TypeWinPerc, winPerLen, -1, winPerPerc, res)
				}()

				//--- Check if we have to stop the process

				if f.stopping {
					slog.Info("generateWinPerc: Got stop request")
					return true
				}
			}
		}
	}

	return false
}

//=============================================================================

func (f *OptimizationProcess) generateEquVsAvg() bool {
	if f.params.EnableEquAvg {
		slog.Info("generateEquVsAvg: Optimizing equity vs its average", "tsId", f.ts.Id, "tsName", f.ts.Name)

		for _, equAvgLen := range f.params.EquAvgLen.Steps() {
			filter := &db.TradingFilter{
				EquAvgEnabled: true,
				EquAvgLen    : equAvgLen,
			}

			go func(){
				res := RunAnalysis(f.ts, filter, f.data)
				f.addResult(TypeEquVsAvg, equAvgLen, -1, -1, res)
			}()

			//--- Check if we have to stop the process

			if f.stopping {
				slog.Info("generateEquVsAvg: Got stop request")
				return true
			}
		}
	}

	return false
}

//=============================================================================

func (f *OptimizationProcess) generateTrendline() bool {
	if f.params.EnableTrendline {
		slog.Info("generateTrendline: Optimizing trendline", "tsId", f.ts.Id, "tsName", f.ts.Name)

		for _, trendLen := range f.params.TrendlineLen.Steps() {
			for _, trendValue := range f.params.TrendlineValue.Steps() {
				filter := &db.TradingFilter{
					TrendlineEnabled: true,
					TrendlineLen    : trendLen,
					TrendlineValue  : trendValue,
				}

				go func(){
					res := RunAnalysis(f.ts, filter, f.data)
					f.addResult(TypeTrendline, trendLen, -1, trendValue, res)
				}()

				//--- Check if we have to stop the process

				if f.stopping {
					slog.Info("generateTrendline: Got stop request")
					return true
				}
			}
		}
	}

	return false
}

//=============================================================================

func (f *OptimizationProcess) addResult(filter string, len int, newLen int, percentage int, far *AnalysisResponse) {
	r := &Run{
		FilterType: filter,
		Length    : len,
		NewLength : newLen,
		Percentage: percentage,
	}

	r.NetProfit   = far.Summary.FilProfit
	r.AvgTrade    = far.Summary.FilAverageTrade
	r.MaxDrawdown = far.Summary.FilMaxDrawdown

	value := f.runComparator.getValue(&far.Summary)
	f.info.addResult(r, value)
}

//=============================================================================

func (f *OptimizationProcess) initValues() {

	//--- Run without filters to get the baseline

	res := RunAnalysis(f.ts, &db.TradingFilter{}, f.data)

	f.info.BaseValue = f.runComparator.getValue(&res.Summary)
	f.info.BestValue = f.info.BaseValue
}

//=============================================================================
