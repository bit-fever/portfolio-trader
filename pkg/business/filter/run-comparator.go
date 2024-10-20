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
	"github.com/emirpasic/gods/utils"
)

//=============================================================================
//===
//=== Run comparator
//===
//=== Notes:
//===  - in reverse order: max to min
//===  - first, compare field to optimize
//=============================================================================

type RunComparator struct {
	compare  utils.Comparator
	getValue func(s *Summary) float64
}

//=============================================================================

func NewRunComparator(field string) *RunComparator {
	if field == FieldToOptimizeNetProfit {
		return &RunComparator{
			compare : compareNetProfit,
			getValue: func(s *Summary) float64{ return s.FilProfit },
		}
	}

	if field == FieldToOptimizeAvgTrade {
		return &RunComparator{
			compare : compareAvgTrade,
			getValue: func(s *Summary) float64{ return s.FilAverageTrade },
		}
	}

	if field == FieldToOptimizeDrawDown {
		return &RunComparator{
			compare : compareDrawDown,
			getValue: func(s *Summary) float64{ return s.FilMaxDrawdown },
		}
	}

	return nil
}

//=============================================================================

func compareNetProfit(a any, b any) int {
	v1 := a.(*Run)
	v2 := b.(*Run)

	if v1.NetProfit < v2.NetProfit { return +1 }
	if v1.NetProfit > v2.NetProfit { return -1 }

	return compareOtherFields(v1, v2)
}

//=============================================================================

func compareAvgTrade(a any, b any) int {
	v1 := a.(*Run)
	v2 := b.(*Run)

	if v1.AvgTrade < v2.AvgTrade { return +1 }
	if v1.AvgTrade > v2.AvgTrade { return -1 }

	return compareOtherFields(v1, v2)
}

//=============================================================================

func compareDrawDown(a any, b any) int {
	v1 := a.(*Run)
	v2 := b.(*Run)

	if v1.MaxDrawdown < v2.MaxDrawdown { return +1 }
	if v1.MaxDrawdown > v2.MaxDrawdown { return -1 }

	return compareOtherFields(v1, v2)
}

//=============================================================================

func compareOtherFields(a,b *Run) int {
	if a.FilterType < b.FilterType { return -1 }
	if a.FilterType > b.FilterType { return +1 }

	if a.Length < b.Length { return -1 }
	if a.Length > b.Length { return +1 }

	if a.NewLength < b.NewLength { return -1 }
	if a.NewLength > b.NewLength { return +1 }

	if a.Percentage < b.Percentage { return -1 }
	if a.Percentage > b.Percentage { return +1 }

	return 0
}

//=============================================================================
