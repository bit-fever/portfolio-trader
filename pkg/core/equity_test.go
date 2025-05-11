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

import (
	"golang.org/x/exp/slices"
	"testing"
)

//=============================================================================

var trades = []float64{ 4, 2, -3, -1, 2 }
var equity = []float64{ 4, 6,  3,  2, 4 }

//=============================================================================

func TestBuildGrossEquity(t *testing.T) {
	equ := BuildEquity(&trades)

	if !slices.Equal(*equ, equity) {
		t.Errorf("Bad gross equity calculation. Expected %v but got %v", equity, equ)
	}
}

//=============================================================================

var equity1 = []float64{ 3, 5,  4, 7, 1 }
var ddown1  = []float64{ 0, 0, -1, 0, -6 }
var maxDd1  = -6.0

//=============================================================================

func TestBuildDrawDown(t *testing.T) {
	dd1, max1 := BuildDrawDown(&equity1)

	if !slices.Equal(*dd1, ddown1) {
		t.Errorf("Bad drawdown result. Got %v", dd1)
	}

	if max1 != maxDd1 {
		t.Errorf("Bad max drawdown result. Expected %v but got %v", maxDd1, max1)
	}
}

//=============================================================================
