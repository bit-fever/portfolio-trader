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

package stats

import (
	"math"
	"testing"
)

//=============================================================================

var prices = []float64{ 1.64, 5.85, 9.22, 3.51, -0.88, 1.07, 13.03, 9.4, 10.49, -5.08, 0, 0 }

//=============================================================================

func TestSharpeRatio(t *testing.T) {
	t.Logf("Avg: %v", Mean(prices))
	t.Logf("Std: %v", StdDev(prices, Mean(prices)))
	t.Logf("SR : %v", SharpeRatio(10, 0.05))
	t.Logf("SRA: %v", SharpeRatio(10, 0.05)*math.Sqrt(12))

	//if !slices.Equal(*equ, equity) {
	//	t.Errorf("Bad gross equity calculation. Expected %v but got %v", equity, equ)
	//}
}

//=============================================================================
