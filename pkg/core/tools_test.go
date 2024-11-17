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

package core

import (
	"golang.org/x/exp/slices"
	"testing"
	"time"
)

//=============================================================================

var days = []time.Time{
	time.Date(2024, 1 , 1, 10, 11, 12, 0, time.UTC),
	time.Date(2024, 2 ,21,  3, 11, 12, 0, time.UTC),
	time.Date(2024, 7 ,11, 14, 11, 12, 0, time.UTC),
	time.Date(2025, 1 ,23, 22, 11, 12, 0, time.UTC),
}

var xAxis = []float64{ 0, 1217, 4612, 9324 }

var yAxis1= []float64{ 125,  87,  90, 130 }
var yAxis2= []float64{ 250, -30, 120, -12 }

//=============================================================================

func TestMean(t *testing.T) {
	mean := calcMean(xAxis)

	if mean != 3788.25 {
		t.Errorf("Bad mean: Expected 3788.25 and got %v", mean)
	}

	mean = calcMean(yAxis1)

	if mean != 108 {
		t.Errorf("Bad mean: Expected 108 and got %v", mean)
	}

	mean = calcMean(yAxis2)

	if mean != 82 {
		t.Errorf("Bad mean: Expected 82 and got %v", mean)
	}
}

//=============================================================================

func TestXAxisCalculation(t *testing.T) {
	axis := calcXAxis(days)

	if !slices.Equal(axis, xAxis) {
		t.Errorf("Bad xAxis result. Got %v", axis)
	}
}

//=============================================================================

func TestLinearRegression(t *testing.T) {
	slope := LinearRegression(days, yAxis1)

	if slope < 0.00184 || slope > 0.00185 {
		t.Errorf("Bad linear regression: Expected ~0.00184 and got %v", slope)
	}

	slope = LinearRegression(days, yAxis2)

	if slope < -0.01602 || slope > -0.01601 {
		t.Errorf("Bad linear regression: Expected ~-0.01601 and got %v", slope)
	}
}

//=============================================================================
