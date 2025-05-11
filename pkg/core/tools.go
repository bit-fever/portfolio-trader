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
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"time"
)

//=============================================================================

func Trunc2d(value float64) float64 {
	return float64(int(value * 100)) / 100
}

//=============================================================================

func LinearRegression(x []time.Time, y []float64) float64 {
	xAxis := calcXAxis(x)
	xMean := calcMean(xAxis)
	yMean := calcMean(y)

	num := 0.0
	den := 0.0
	aux := 0.0

	for i,_ := range y {
		aux = xAxis[i] - xMean
		num += aux * (y[i] - yMean)
		den += aux * aux
	}

	return num / den
}

//=============================================================================

func GetLocation(timezone string, ts *db.TradingSystem) (*time.Location, error) {
	if timezone == "exchange" {
		timezone = ts.Timezone
	}

	return time.LoadLocation(timezone)
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func calcXAxis(time []time.Time) []float64 {
	var x []float64

	for _, t := range time {
		v := calcHours(t, time[0])
		x = append(x, v)
	}

	return x
}

//=============================================================================

func calcHours(t,t0 time.Time) float64 {
	delta := t.UnixMilli() - t0.UnixMilli()
	return float64(delta / 1000 / 3600)
}

//=============================================================================

func calcMean(data []float64) float64 {
	sum := 0.0

	for _,value := range data {
		sum += value
	}

	return sum/float64(len(data))
}

//=============================================================================
