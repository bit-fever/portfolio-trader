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

package core

//=============================================================================

func CalcDrawDown(equity *[]float64, drawDown *[]float64) float64 {
	maxProfit    := 0.0
	currDrawDown := 0.0
	maxDrawDown  := 0.0

	for i, currProfit := range *equity {
		if currProfit >= maxProfit {
			maxProfit = currProfit
			currDrawDown = 0
		} else {
			currDrawDown = currProfit - maxProfit
		}

		(*drawDown)[i] = currDrawDown

		if currDrawDown < maxDrawDown {
			maxDrawDown = currDrawDown
		}
	}

	return maxDrawDown
}

//=============================================================================

func CalcWinningPercentage(profits []float64, filter []int8) float64 {
	tot := 0
	pos := 0

	for i, profit := range profits {
		if profit != 0 {
			if filter == nil || filter[i] == 1 {
				tot++
				if profit > 0 {
					pos++
				}
			}
		}
	}

	if tot == 0 {
		return 0
	}

	return float64(pos * 10000 / tot) / 100
}

//=============================================================================

func CalcAverageTrade(profits []float64, filter []int8) float64 {
	sum := 0.0
	num := 0.0

	for i, profit := range profits {
		if profit != 0 {
			if filter == nil || filter[i] == 1 {
				sum += profit
				num++
			}
		}
	}

	return float64(int(sum * 100 / num)) / 100
}

//=============================================================================
