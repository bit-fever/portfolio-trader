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

package simple

import (
	"github.com/bit-fever/portfolio-trader/pkg/business/filter/algorithm/optimization"
)

//=============================================================================

type simpleAlgorithm struct {
	ctx  optimization.Context
	fc  *optimization.FilterConfig
}

//=============================================================================

func New() optimization.Algorithm {
	return &simpleAlgorithm{}
}

//=============================================================================
//===
//=== Simple algorithm implementation
//===
//=============================================================================

func (sa *simpleAlgorithm) Init(ctx optimization.Context) {
	sa.ctx = ctx
	sa.fc  = ctx.FilterConfig()
}

//=============================================================================

func (sa *simpleAlgorithm) StepsCount() uint {
	return	sa.stepsCountPosProfit() +
			sa.stepsCountOldNew   () +
			sa.stepsCountWinPerc  () +
			sa.stepsCountEquAvg   () +
			sa.stepsCountTrendline()
}

//=============================================================================

func (sa *simpleAlgorithm) Optimize() {
	if !sa.generatePosProfit(){
		if !sa.generateOldVsNew() {
			if !sa.generateWinPerc() {
				if !sa.generateEquVsAvg() {
					sa.generateTrendline()
				}
			}
		}
	}
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func (sa *simpleAlgorithm) stepsCountPosProfit() uint {
	if sa.fc.EnablePosProfit {
		return sa.fc.PosProLen.StepsCount()
	}

	return 0
}

//=============================================================================

func (sa *simpleAlgorithm) stepsCountOldNew() uint {
	if sa.fc.EnableOldNew {
		return sa.fc.OldNewOldLen.StepsCount() * sa.fc.OldNewNewLen.StepsCount() * sa.fc.OldNewOldPerc.StepsCount()
	}

	return 0
}

//=============================================================================

func (sa *simpleAlgorithm) stepsCountWinPerc() uint {
	if sa.fc.EnableWinPerc {
		return sa.fc.WinPercLen.StepsCount() * sa.fc.WinPercPerc.StepsCount()
	}

	return 0
}

//=============================================================================

func (sa *simpleAlgorithm) stepsCountEquAvg() uint {
	if sa.fc.EnableEquAvg {
		return sa.fc.EquAvgLen.StepsCount()
	}

	return 0
}

//=============================================================================

func (sa *simpleAlgorithm) stepsCountTrendline() uint {
	if sa.fc.EnableTrendline {
		return sa.fc.TrendlineLen.StepsCount() * sa.fc.TrendlineValue.StepsCount()
	}

	return 0
}

//=============================================================================

func (sa *simpleAlgorithm) generatePosProfit() bool {
	if sa.fc.EnablePosProfit {
		sa.ctx.LogInfo("generatePosProfit: Optimizing positive profit")

		for _, posProLen := range *sa.fc.PosProLen.Steps() {
			f := sa.ctx.Baseline()
			f.PosProEnabled= true
			f.PosProLen    = posProLen

			go func() {
				sa.ctx.RunAnalysis(&f)
			}()

			//--- Check if we have to stop the process

			if sa.ctx.IsStopping() {
				sa.ctx.LogInfo("generatePosProfit: Got stop request")
				return true
			}
		}
	}

	return false
}

//=============================================================================

func (sa *simpleAlgorithm) generateOldVsNew() bool {
	if sa.fc.EnableOldNew {
		sa.ctx.LogInfo("generateOldVsNew: Optimizing old vs new periods")

		for _, oldNewOldLen := range *sa.fc.OldNewOldLen.Steps() {
			for _, oldNewNewLen := range *sa.fc.OldNewNewLen.Steps() {
				for _, oldNewOldPerc := range *sa.fc.OldNewOldPerc.Steps() {
					f := sa.ctx.Baseline()
					f.OldNewEnabled = true
					f.OldNewOldLen  = oldNewOldLen
					f.OldNewNewLen  = oldNewNewLen
					f.OldNewOldPerc = oldNewOldPerc

					go func() {
						sa.ctx.RunAnalysis(&f)
					}()

					//--- Check if we have to stop the process

					if sa.ctx.IsStopping() {
						sa.ctx.LogInfo("generateOldVsNew: Got stop request")
						return true
					}
				}
			}
		}
	}

	return false
}

//=============================================================================

func (sa *simpleAlgorithm) generateWinPerc() bool {
	if sa.fc.EnableWinPerc {
		sa.ctx.LogInfo("generateWinPerc: Optimizing winning percentage")

		for _, winPerLen := range *sa.fc.WinPercLen.Steps() {
			for _, winPerPerc := range *sa.fc.WinPercPerc.Steps() {
				f := sa.ctx.Baseline()
				f.WinPerEnabled= true
				f.WinPerLen    = winPerLen
				f.WinPerValue  = winPerPerc

				go func(){
					sa.ctx.RunAnalysis(&f)
				}()

				//--- Check if we have to stop the process

				if sa.ctx.IsStopping() {
					sa.ctx.LogInfo("generateWinPerc: Got stop request")
					return true
				}
			}
		}
	}

	return false
}

//=============================================================================

func (sa *simpleAlgorithm) generateEquVsAvg() bool {
	if sa.fc.EnableEquAvg {
		sa.ctx.LogInfo("generateEquVsAvg: Optimizing equity vs its average")

		for _, equAvgLen := range *sa.fc.EquAvgLen.Steps() {
			f := sa.ctx.Baseline()
			f.EquAvgEnabled= true
			f.EquAvgLen    = equAvgLen

			go func(){
				sa.ctx.RunAnalysis(&f)
			}()

			//--- Check if we have to stop the process

			if sa.ctx.IsStopping() {
				sa.ctx.LogInfo("generateEquVsAvg: Got stop request")
				return true
			}
		}
	}

	return false
}

//=============================================================================

func (sa *simpleAlgorithm) generateTrendline() bool {
	if sa.fc.EnableTrendline {
		sa.ctx.LogInfo("generateTrendline: Optimizing trendline")

		for _, trendLen := range *sa.fc.TrendlineLen.Steps() {
			for _, trendValue := range *sa.fc.TrendlineValue.Steps() {
				f := sa.ctx.Baseline()
				f.TrendlineEnabled= true
				f.TrendlineLen    = trendLen
				f.TrendlineValue  = trendValue

				go func(){
					sa.ctx.RunAnalysis(&f)
				}()

				//--- Check if we have to stop the process

				if sa.ctx.IsStopping() {
					sa.ctx.LogInfo("generateTrendline: Got stop request")
					return true
				}
			}
		}
	}

	return false
}

//=============================================================================
