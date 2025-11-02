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
	"github.com/tradalia/portfolio-trader/pkg/core"
	"math"
)

//=============================================================================
//===
//=== Run comparator
//===
//=== Notes:
//===  - in reverse order: max to min
//=============================================================================

func runComparator(a any, b any) int {
	r1 := a.(*Run)
	r2 := b.(*Run)
	v1 := r1.FitnessValue
	v2 := r2.FitnessValue

	if v1 < v2 { return +1 }
	if v1 > v2 { return -1 }

	if r1.random == r2.random {
		return 0
	}

	if r1.random < r2.random {
		return 1
	}

	return -1
}

//=============================================================================
//===
//=== Fitness functions
//===
//=============================================================================

type FitnessFunction func(r *Run) float64

//=============================================================================

func GetFitnessFunction(field string) FitnessFunction {
	switch field {
		case FieldToOptimizeNetProfit:
			return ffNetProfit

		case FieldToOptimizeAvgTrade:
			return ffAvgTrade

		case FieldToOptimizeDrawDown:
			return ffMaxDrawdown

		case FieldToOptimizeNetProfitAvgTrade:
			return ffNetProfitAvgTrade

		case FieldToOptimizeNetProfitAvgTradeMaxDD:
			return ffNetProfitAvgTradeMaxDD

		default:
			panic("Unknown field to optimize: "+ field)
	}
}

//=============================================================================

func ffNetProfit(r *Run) float64 {
	return r.NetProfit
}

//=============================================================================

func ffAvgTrade(r *Run) float64 {
	return r.AvgTrade
}

//=============================================================================

func ffMaxDrawdown(r *Run) float64 {
	return r.MaxDrawdown
}

//=============================================================================

func ffNetProfitAvgTrade(r *Run) float64 {
	fv := r.NetProfit*r.AvgTrade

	//--- We have to push down results where both netProfit and avgTrade are negative
	//--- because their product is positive

	if r.NetProfit < 0 && r.AvgTrade < 0 {
		fv *= -1
	}

	if fv <= -1000 || fv >= 1000 {
		fv = math.Round(fv)
	}

	return fv
}

//=============================================================================

func ffNetProfitAvgTradeMaxDD(r *Run) float64 {
	fv := r.NetProfit*r.AvgTrade
	dd := math.Abs(r.MaxDrawdown)

	//--- We have to push down results where both netProfit and avgTrade are negative
	//--- because their product is positive

	if r.NetProfit < 0 && r.AvgTrade < 0 {
		fv *= -1
	}

	if dd == 0 {
		dd = 1
	}

	fv = core.Trunc2d(fv/dd)

	if fv <= -1000 || fv >= 1000 {
		fv = math.Round(fv)
	}

	return fv
}

//=============================================================================
