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

package filter

import (
	"github.com/bit-fever/portfolio-trader/pkg/core"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"time"
)

//=============================================================================
//===
//=== Filter activation calculation
//===
//=============================================================================

func CalcActivation(ts *db.TradingSystem, filter *db.TradingFilter, list []db.Trade) bool {
	if len(list) == 0 {
		return true
	}

	e := &Equities{}

	//--- Calc unfiltered equity and days
	calcUnfilteredEquityAndProfit(e, ts, &list)

	if filter.EquAvgEnabled {
		e.Average = calcAverageEquity(e.Time, e.UnfilteredEquity, filter.EquAvgLen)
	}

	a := calcActivations(e, filter)

	return a.IsLastActive()
}

//=============================================================================
//===
//=== AnalysisResponse building
//===
//=============================================================================

func RunAnalysis(ts *db.TradingSystem, filter *db.TradingFilter, list *[]db.Trade) *AnalysisResponse {
	res := &AnalysisResponse{}
	res.TradingSystem.Id   = ts.Id
	res.TradingSystem.Name = ts.Name
	res.Filter             = filter

	//--- Creates slices

	if len(*list) == 0 {
		return res
	}

	e := &res.Equities

	//--- Calc unfiltered equity and days
	calcUnfilteredEquityAndProfit(e, ts, list)

	if res.Filter.EquAvgEnabled {
		e.Average = calcAverageEquity(e.Time, e.UnfilteredEquity, res.Filter.EquAvgLen)
	}

	res.Activations = calcActivations(e, filter)
	calcFilterActivation(e, res.Activations, filter)
	calcFilteredEquity(res)

	unfilteredDrawdown, maxUnfDD := core.BuildDrawDown(&e.UnfilteredEquity)
	filteredDrawdown,   maxFilDD := core.BuildDrawDown(&e.FilteredEquity)

	e.UnfilteredDrawdown = *unfilteredDrawdown
	e.FilteredDrawdown   = *filteredDrawdown

	calcSummary(res, maxUnfDD, maxFilDD)

	return res
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func calcUnfilteredEquityAndProfit(e *Equities, ts *db.TradingSystem, tradeList *[]db.Trade) {
	netEquity     := 0.0
	costPerOperat := float64(ts.CostPerOperation)
	size          := len(*tradeList)

	e.Time              = make([]time.Time,size)
	e.NetProfit         = make([]float64,  size)
	e.UnfilteredEquity  = make([]float64,  size)

	for i, t := range *tradeList {
		netProfit := t.GrossProfit - costPerOperat * 2
		netEquity += netProfit

		e.Time[i]             = *t.ExitDate
		e.UnfilteredEquity[i] = netEquity
		e.NetProfit[i]        = netProfit
	}
}

//=============================================================================

func calcAverageEquity(times []time.Time, equity []float64, maLen int) *core.Serie {
	p := core.Serie{}

	maSum := 0.0

	for i, t := range times {
		maSum += equity[i]

		if i >= maLen -1 {
			if i>maLen -1 {
				maSum -= equity[i-maLen]
			}
			p.AddPoint(t, maSum/float64(maLen))
		} else {
			p.AddPoint(t, equity[i])
		}
	}

	//--- If we can't calculate the average (days is too high), just return nil
	if p.Time == nil {
		return nil
	}

	return &p
}

//=============================================================================
//===
//=== Calculate activations
//===
//=============================================================================

func calcActivations(e *Equities, f *db.TradingFilter) *Activations {
	a := &Activations{}
	a.EquityVsAverage   = calcEquAvgActivation   (e, f)
	a.PositiveProfit    = calcPosProfitActivation(e, f)
	a.WinningPercentage = calcWinPercActivation  (e, f)
	a.OldVsNew          = calcOldVsNewActivation (e, f)
	a.Trendline         = calcTrendlineActivation(e, f)
	a.Drawdown          = calcDrawdownActivation (e, f)
	return a
}

//=============================================================================

func calcEquAvgActivation(e *Equities, f *db.TradingFilter) *Activation {
	if !f.EquAvgEnabled || e.Average == nil {
		return nil
	}

	a := Activation{}

	avg := e.Average

	for i, avgTime := range avg.Time {
		if i == 0 {
			a.AddPoint(avgTime, 1)
		} else {
			avgVal := avg.Values[i]
			equVal := e.UnfilteredEquity[i]
			value  := int8(0)

			if equVal >= avgVal {
				value = 1
			}

			a.AddPoint(avgTime, value)
		}
	}

	return &a
}

//=============================================================================

func calcPosProfitActivation(e *Equities, f *db.TradingFilter) *Activation {
	if !f.PosProEnabled {
		return nil
	}

	a := Activation{}

	profSum  := 0.0
	profDays := f.PosProLen
	equity   := e.UnfilteredEquity

	for i, t := range e.Time {
		if i >= profDays -1 {
			profSum = equity[i]

			if i>profDays -1 {
				profSum -= equity[i-profDays]
			}

			value := int8(0)

			if profSum >= 0 {
				value = 1
			}
			a.AddPoint(t, value)
		}
	}

	//--- If we can't calculate the average (days is too high), just return nil
	if a.Time == nil {
		return nil
	}

	return &a
}

//=============================================================================

func calcWinPercActivation(e *Equities, f *db.TradingFilter) *Activation {
	if !f.WinPerEnabled {
		return nil
	}

	a := Activation{}

	posCount := 0
	totCount := 0
	winLen   := f.WinPerLen
	percValue:= f.WinPerValue
	profits  := e.NetProfit

	for i, t := range e.Time {
		if profits[i] != 0 {
			totCount++

			if profits[i] > 0 {
				posCount++
			}
		}

		if i >= winLen -1 {
			if i>winLen -1 {
				if profits[i-winLen] != 0 {
					totCount--

					if profits[i-winLen] > 0 {
						posCount--
					}
				}
			}

			value := int8(0)
			if totCount > 0 {
				if posCount * 100 / totCount >= percValue {
					value = 1
				}
			}
			a.AddPoint(t, value)
		}
	}

	//--- If we can't calculate the average (days is too high), just return nil
	if a.Time == nil {
		return nil
	}

	return &a
}

//=============================================================================

func calcOldVsNewActivation(e *Equities, f *db.TradingFilter) *Activation {
	if !f.OldNewEnabled {
		return nil
	}

	a := Activation{}

	oldSum  := 0.0
	newSum  := 0.0
	oldLen  := f.OldNewOldLen
	newLen  := f.OldNewNewLen
	equity  := e.UnfilteredEquity
	oldPerc := float64(f.OldNewOldPerc)/100.0

	for i, t := range e.Time {
		//--- New period

		if i >= newLen -1 {
			newSum = equity[i]

			if i>newLen -1 {
				newSum -= equity[i-newLen]
			}

			//--- Old period

			if i >= (oldLen+newLen) -1 {
				oldSum = equity[i-newLen]

				if i>(oldLen+newLen) -1 {
					oldSum -= equity[i-newLen-oldLen]
				}

				value := int8(0)

				if newSum >= oldSum * oldPerc {
					value = 1
				}
				a.AddPoint(t, value)
			}
		}
	}

	//--- If we can't calculate the average (days is too high), just return nil
	if a.Time == nil {
		return nil
	}

	return &a
}

//=============================================================================

func calcTrendlineActivation(e *Equities, f *db.TradingFilter) *Activation {
	if !f.TrendlineEnabled {
		return nil
	}

	a := Activation{}

	trendLen:= f.TrendlineLen
	thresh  := float64(f.TrendlineValue) / 100
	equity  := e.UnfilteredEquity

	for i, t := range e.Time {
		if i >= trendLen -1 {
			slope := core.LinearRegression(e.Time[i -trendLen +1:i], equity[i -trendLen +1:i])
			value := int8(0)

			if slope >= thresh {
				value = 1
			}
			a.AddPoint(t, value)
		}
	}

	//--- If we can't calculate the average (days is too high), just return nil
	if a.Time == nil {
		return nil
	}

	return &a
}

//=============================================================================

func calcDrawdownActivation(e *Equities, f *db.TradingFilter) *Activation {
	if !f.DrawdownEnabled {
		return nil
	}

	a := Activation{}

	minDD        := float64(f.DrawdownMin)
	maxDD        := float64(f.DrawdownMax)
	maxProfit    := 0.0
	currDrawDown := 0.0
	value        := int8(1)

	for i, currProfit := range e.UnfilteredEquity {
		if currProfit >= maxProfit {
			maxProfit = currProfit
			currDrawDown = 0
		} else {
			currDrawDown = currProfit - maxProfit
		}

		if currDrawDown < -maxDD {
			value = 0
		} else if currDrawDown > -minDD {
			value = 1
		}

		a.AddPoint(e.Time[i], value)
	}

	return &a
}

//=============================================================================

func calcFilterActivation(e *Equities, a*Activations, f *db.TradingFilter) {
	avgEquStrategy   := NewActivationStrategy(a.EquityVsAverage,   f.EquAvgEnabled)
	posProfStrategy  := NewActivationStrategy(a.PositiveProfit,    f.PosProEnabled)
	winPerStrategy   := NewActivationStrategy(a.WinningPercentage, f.WinPerEnabled)
	oldNewStrategy   := NewActivationStrategy(a.OldVsNew,          f.OldNewEnabled)
	trendStrategy    := NewActivationStrategy(a.Trendline,         f.TrendlineEnabled)
	drawdownStrategy := NewActivationStrategy(a.Drawdown,          f.DrawdownEnabled)

	e.FilterActivation = make([]int8, len(e.Time))

	for i, t := range e.Time {
		//--- These 6 conditions must be standalone. If we use tot := A && B && C && D
		//--- then B, C, D evaluation can be skipped because the && operator is already satisfied

		avgEqu := avgEquStrategy  .IsActive(t)
		posPro := posProfStrategy .IsActive(t)
		winPer := winPerStrategy  .IsActive(t)
		oldNew := oldNewStrategy  .IsActive(t)
		trend  := trendStrategy   .IsActive(t)
		drawd  := drawdownStrategy.IsActive(t)

		if avgEqu && posPro && winPer && oldNew && trend && drawd {
			e.FilterActivation[i] = 1
		}
	}
}

//=============================================================================

func calcFilteredEquity(res *AnalysisResponse) {
	equ := &res.Equities
	sum := 0.0

	equ.FilteredEquity = make([]float64, len(equ.NetProfit))

	for i, value := range equ.NetProfit {
		if i>0 {
			//--- We have to use the activation at time [i-1]
			if equ.FilterActivation[i-1] == 0 {
				value = float64(0)
			}
		}

		sum += value
		equ.FilteredEquity[i] = sum
	}
}

//=============================================================================

func calcSummary(res *AnalysisResponse, maxUnfDD, maxFilDD float64) {
	sum  := &res.Summary
	last := len(res.Equities.Time) -1

	sum.UnfProfit      = res.Equities.UnfilteredEquity[last]
	sum.FilProfit      = res.Equities.FilteredEquity[last]
	sum.UnfMaxDrawdown = maxUnfDD
	sum.FilMaxDrawdown = maxFilDD
	sum.UnfWinningPerc = core.CalcWinningPercentage(res.Equities.NetProfit, nil)
	sum.FilWinningPerc = core.CalcWinningPercentage(res.Equities.NetProfit, res.Equities.FilterActivation)
	sum.UnfAverageTrade= core.CalcAverageTrade     (res.Equities.NetProfit, nil)
	sum.FilAverageTrade= core.CalcAverageTrade     (res.Equities.NetProfit, res.Equities.FilterActivation)
}

//=============================================================================
