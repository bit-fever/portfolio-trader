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
//=== AnalysisResponse building
//===
//=============================================================================

func RunAnalysis(ts *db.TradingSystem, filters *TradingFilters, list *[]db.Trade) *AnalysisResponse {
	res := &AnalysisResponse{}
	res.TradingSystem.Id   = ts.Id
	res.TradingSystem.Name = ts.Name
	res.Filters            = *filters

	//--- Creates slices

	size := len(*list)

	if size == 0 {
		return res
	}

	e := &res.Equities
	e.Time              = make([]time.Time,size)
	e.NetProfit         = make([]float64,  size)
	e.UnfilteredEquity  = make([]float64,  size)
	e.FilteredEquity    = make([]float64,  size)
	e.UnfilteredDrawdown= make([]float64,  size)
	e.FilteredDrawdown  = make([]float64,  size)
	e.FilterActivation  = make([]int8,     size)

	//--- Calc unfiltered equity and days
	calcUnfilteredEquityAndProfit(e, ts, list)

	if res.Filters.EquAvgEnabled {
		e.Average = calcAverageEquity(e.Time, e.UnfilteredEquity, res.Filters.EquAvgLen)
	}

	calcActivations(res)
	calcFilterActivation(res)
	calcFilteredEquity(res)

	maxUnfDD := core.CalcDrawDown(&e.UnfilteredEquity, &e.UnfilteredDrawdown)
	maxFilDD := core.CalcDrawDown(&e.FilteredEquity,   &e.FilteredDrawdown)

	calcSummary(res, maxUnfDD, maxFilDD)

	return res
}

//=============================================================================

func calcUnfilteredEquityAndProfit(e *Equities, ts *db.TradingSystem, tradeList *[]db.Trade) {
	currProfit    := 0.0
	costPerOperat := float64(ts.CostPerOperation)

	for i, t := range *tradeList {
		currCost   := t.GrossProfit - costPerOperat * 2
		currProfit += currCost

		e.Time[i]             = *t.ExitTime
		e.UnfilteredEquity[i] = currProfit
		e.NetProfit[i]        = currCost
	}
}

//=============================================================================
//===
//=== Equity average filtering
//===
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

func calcActivations(res *AnalysisResponse) {
	res.Activations.EquityVsAverage   = calcEquAvgActivation   (res)
	res.Activations.PositiveProfit    = calcPosProfitActivation(res)
	res.Activations.WinningPercentage = calcWinPercActivation  (res)
	res.Activations.OldVsNew          = calcOldVsNewActivation (res)
}

//=============================================================================

func calcEquAvgActivation(res *AnalysisResponse) *Activation {
	if !res.Filters.EquAvgEnabled || res.Equities.Average == nil {
		return nil
	}

	a := Activation{}

	avg     := res.Equities.Average
	avgDays := res.Filters.EquAvgLen

	for i, avgTime := range avg.Time {
		avgVal := avg.Values[i]
		equVal := res.Equities.UnfilteredEquity[i + avgDays -1]
		value  := int8(0)

		if equVal >= avgVal {
			value = 1
		}

		a.AddPoint(avgTime, value)
	}

	return &a
}

//=============================================================================

func calcPosProfitActivation(res *AnalysisResponse) *Activation {
	if !res.Filters.PosProEnabled {
		return nil
	}

	a := Activation{}

	profSum  := 0.0
	profDays := res.Filters.PosProLen
	equity   := res.Equities.UnfilteredEquity

	for i, t := range res.Equities.Time {
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

func calcWinPercActivation(res *AnalysisResponse) *Activation {
	if !res.Filters.WinPerEnabled {
		return nil
	}

	a := Activation{}

	posCount := 0
	totCount := 0
	winLen   := res.Filters.WinPerLen
	percValue:= res.Filters.WinPerValue
	profits  := res.Equities.NetProfit

	for i, t := range res.Equities.Time {
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

func calcOldVsNewActivation(res *AnalysisResponse) *Activation {
	if !res.Filters.OldNewEnabled {
		return nil
	}

	a := Activation{}

	oldSum  := 0.0
	newSum  := 0.0
	oldLen  := res.Filters.OldNewOldLen
	newLen  := res.Filters.OldNewNewLen
	equity  := res.Equities.UnfilteredEquity
	oldPerc := float64(res.Filters.OldNewOldPerc)/100.0

	for i, t := range res.Equities.Time {
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

func calcFilterActivation(res *AnalysisResponse) {
	equ := &res.Equities
	act := &res.Activations
	fil := &res.Filters

	avgEquStrategy  := NewActivationStrategy(act.EquityVsAverage,   fil.EquAvgEnabled)
	posProfStrategy := NewActivationStrategy(act.PositiveProfit,    fil.PosProEnabled)
	winPerStrategy  := NewActivationStrategy(act.WinningPercentage, fil.WinPerEnabled)
	oldNewStrategy  := NewActivationStrategy(act.OldVsNew,          fil.OldNewEnabled)

	for i, t := range equ.Time {
		//--- These 4 conditions must be standalone. If we use tot := A && B && C && D
		//--- then B, C, D evaluation can be skipped because the && operator is already satisfied

		avgEqu := avgEquStrategy .IsActive(t)
		posPro := posProfStrategy.IsActive(t)
		winPer := winPerStrategy .IsActive(t)
		oldNew := oldNewStrategy .IsActive(t)

		if avgEqu && posPro && winPer && oldNew {
			equ.FilterActivation[i] = 1
		}
	}
}

//=============================================================================

func calcFilteredEquity(res *AnalysisResponse) {
	equ := &res.Equities
	sum := 0.0

	for i, value := range equ.NetProfit {
		if equ.FilterActivation[i] == 0 {
			value = float64(0)
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
	sum.UnfAverageTrade= core.CalcAverageTrade(res.Equities.NetProfit, nil)
	sum.FilAverageTrade= core.CalcAverageTrade(res.Equities.NetProfit, res.Equities.FilterActivation)
}

//=============================================================================
