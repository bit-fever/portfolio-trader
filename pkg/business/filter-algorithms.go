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
	"github.com/bit-fever/portfolio-trader/pkg/business/filter"
	"github.com/bit-fever/portfolio-trader/pkg/core"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"github.com/bit-fever/portfolio-trader/pkg/model"
)

//=============================================================================
//===
//=== FilterAnalysisResponse building
//===
//=============================================================================

func buildFilterResponse(ts *db.TradingSystem, filters *TradingFilters, list *[]db.DailyInfo) *FilterAnalysisResponse {
	res := &FilterAnalysisResponse{}
	res.TradingSystem.Id   = ts.Id
	res.TradingSystem.Name = ts.Name
	res.Filters            = *filters

	//--- Creates slices

	size := len(*list)

	if size == 0 {
		return res
	}

	e := &res.Equities
	e.Days              = make([]int,     size)
	e.NetProfits        = make([]float64, size)
	e.UnfilteredEquity  = make([]float64, size)
	e.FilteredEquity    = make([]float64, size)
	e.UnfilteredDrawdown= make([]float64, size)
	e.FilteredDrawdown  = make([]float64, size)
	e.FilterActivation  = make([]int8,    size)

	//--- Calc unfiltered equity and days
	calcUnfilteredEquityAndProfit(e, ts, list)

	if res.Filters.EquAvgEnabled {
		e.Average = calcAverageEquity(e.Days, e.UnfilteredEquity, res.Filters.EquAvgDays)
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

func calcUnfilteredEquityAndProfit(e *Equities, ts *db.TradingSystem, diList *[]db.DailyInfo) {
	currProfit   := 0.0
	costPerTrade := float64(ts.CostPerTrade)

	for i, di := range *diList {
		currCost   := di.ClosedProfit - costPerTrade * float64(di.NumTrades)
		currProfit += currCost

		e.Days[i]             = di.Day
		e.UnfilteredEquity[i] = currProfit
		e.NetProfits[i]       = currCost
	}
}

//=============================================================================
//===
//=== Equity average filtering
//===
//=============================================================================

func calcAverageEquity(days []int, equity []float64, maDays int) *model.Plot {
	p := model.Plot{}

	maSum  := 0.0

	for i, day := range days {
		maSum += equity[i]

		if i >= maDays -1 {
			if i>maDays -1 {
				maSum -= equity[i-maDays]
			}
			p.AddPoint(day, maSum/float64(maDays))
		}
	}

	//--- If we can't calculate the average (days is too high), just return nil
	if p.Days == nil {
		return nil
	}

	return &p
}

//=============================================================================
//===
//=== Calculate activations
//===
//=============================================================================

func calcActivations(res *FilterAnalysisResponse) {
	res.Activations.EquityVsAverage   = calcEquAvgActivation   (res)
	res.Activations.PositiveProfit    = calcPosProfitActivation(res)
	res.Activations.WinningPercentage = calcWinPercActivation  (res)
	res.Activations.OldVsNew          = calcOldVsNewActivation (res)
}

//=============================================================================

func calcEquAvgActivation(res *FilterAnalysisResponse) *filter.Activation {
	if !res.Filters.EquAvgEnabled || res.Equities.Average == nil {
		return nil
	}

	a := filter.Activation{}

	avg     := res.Equities.Average
	avgDays := res.Filters.EquAvgDays

	for i, avgDay := range avg.Days {
		avgVal := avg.Values[i]
		equVal := res.Equities.UnfilteredEquity[i + avgDays -1]
		value  := int8(0)

		if equVal >= avgVal {
			value = 1
		}

		a.AddDay(avgDay, value)
	}

	return &a
}

//=============================================================================

func calcPosProfitActivation(res *FilterAnalysisResponse) *filter.Activation {
	if !res.Filters.PosProEnabled {
		return nil
	}

	a := filter.Activation{}

	profSum  := 0.0
	profDays := res.Filters.PosProDays
	equity   := res.Equities.UnfilteredEquity

	for i, day := range res.Equities.Days {
		if i >= profDays -1 {
			profSum = equity[i]

			if i>profDays -1 {
				profSum -= equity[i-profDays]
			}

			value := int8(0)

			if profSum >= 0 {
				value = 1
			}
			a.AddDay(day, value)
		}
	}

	//--- If we can't calculate the average (days is too high), just return nil
	if a.Days == nil {
		return nil
	}

	return &a
}

//=============================================================================

func calcWinPercActivation(res *FilterAnalysisResponse) *filter.Activation {
	if !res.Filters.WinPerEnabled {
		return nil
	}

	a := filter.Activation{}

	posCount := 0
	totCount := 0
	winDays  := res.Filters.WinPerDays
	percValue:= res.Filters.WinPerValue
	profits  := res.Equities.NetProfits

	for i, day := range res.Equities.Days {
		if profits[i] != 0 {
			totCount++

			if profits[i] > 0 {
				posCount++
			}
		}

		if i >= winDays -1 {
			if i>winDays -1 {
				if profits[i-winDays] != 0 {
					totCount--

					if profits[i-winDays] > 0 {
						posCount--
					}
				}
			}

			value := int8(0)
			if posCount * 100 / totCount >= percValue {
				value = 1
			}
			a.AddDay(day, value)
		}
	}

	//--- If we can't calculate the average (days is too high), just return nil
	if a.Days == nil {
		return nil
	}

	return &a
}

//=============================================================================

func calcOldVsNewActivation(res *FilterAnalysisResponse) *filter.Activation {
	if !res.Filters.OldNewEnabled {
		return nil
	}

	a := filter.Activation{}

	oldSum  := 0.0
	newSum  := 0.0
	oldDays := res.Filters.OldNewOldDays
	newDays := res.Filters.OldNewNewDays
	equity  := res.Equities.UnfilteredEquity
	oldPerc := float64(res.Filters.OldNewOldPerc)/100.0

	for i, day := range res.Equities.Days {
		//--- New period

		if i >= newDays -1 {
			newSum = equity[i]

			if i>newDays -1 {
				newSum -= equity[i-newDays]
			}

			//--- Old period

			if i >= (oldDays+newDays) -1 {
				oldSum = equity[i-newDays]

				if i>(oldDays+newDays) -1 {
					oldSum -= equity[i-newDays-oldDays]
				}

				value := int8(0)

				if newSum >= oldSum * oldPerc {
					value = 1
				}
				a.AddDay(day, value)
			}
		}
	}

	//--- If we can't calculate the average (days is too high), just return nil
	if a.Days == nil {
		return nil
	}

	return &a
}

//=============================================================================

func calcFilterActivation(res *FilterAnalysisResponse) {
	equ := &res.Equities
	act := &res.Activations
	fil := &res.Filters

	avgEquStrategy  := filter.NewActivationStrategy(act.EquityVsAverage,   fil.EquAvgEnabled)
	posProfStrategy := filter.NewActivationStrategy(act.PositiveProfit,    fil.PosProEnabled)
	winPerStrategy  := filter.NewActivationStrategy(act.WinningPercentage, fil.WinPerEnabled)
	oldNewStrategy  := filter.NewActivationStrategy(act.OldVsNew,          fil.OldNewEnabled)

	for i, day := range equ.Days {
		//--- These 4 conditions must be standalone. If we use tot := A && B && C && D
		//--- then B, C, D evaluation can be skipped because the && operator is already satisfied

		avgEqu := avgEquStrategy .IsActive(day)
		posPro := posProfStrategy.IsActive(day)
		winPer := winPerStrategy .IsActive(day)
		oldNew := oldNewStrategy .IsActive(day)

		if avgEqu && posPro && winPer && oldNew {
			equ.FilterActivation[i] = 1
		}
	}
}

//=============================================================================

func calcFilteredEquity(res *FilterAnalysisResponse) {
	equ := &res.Equities
	sum := 0.0

	for i, value := range equ.NetProfits {
		if equ.FilterActivation[i] == 0 {
			value = float64(0)
		}

		sum += value
		equ.FilteredEquity[i] = sum
	}
}

//=============================================================================

func calcSummary(res *FilterAnalysisResponse, maxUnfDD, maxFilDD float64) {
	sum  := &res.Summary
	last := len(res.Equities.Days) -1

	sum.UnfProfit      = res.Equities.UnfilteredEquity[last]
	sum.FilProfit      = res.Equities.FilteredEquity[last]
	sum.UnfMaxDrawdown = maxUnfDD
	sum.FilMaxDrawdown = maxFilDD
	sum.UnfWinningPerc = core.CalcWinningPercentage(res.Equities.NetProfits, nil)
	sum.FilWinningPerc = core.CalcWinningPercentage(res.Equities.NetProfits, res.Equities.FilterActivation)
	sum.UnfAverageTrade= core.CalcAverageTrade(res.Equities.NetProfits, nil)
	sum.FilAverageTrade= core.CalcAverageTrade(res.Equities.NetProfits, res.Equities.FilterActivation)
}

//=============================================================================
