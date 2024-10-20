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

package business

import "time"

//=============================================================================
//===
//=== PortfolioMonitoringResponse
//===
//=============================================================================

type PortfolioMonitoringParams struct {
	TsIds  []uint `form:"tsIds"  binding:"required,min=1,dive"`
	Period    int `form:"period" binding:"required,min=1,max=5000"`
}

//=============================================================================

type BaseMonitoring struct {
	Time          []time.Time `json:"time"`
	GrossProfit   []float64   `json:"grossProfit"`
	NetProfit     []float64   `json:"netProfit"`
	GrossDrawdown []float64   `json:"grossDrawdown"`
	NetDrawdown   []float64   `json:"netDrawdown"`
}

//=============================================================================

type TradingSystemMonitoring struct {
	BaseMonitoring
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

//=============================================================================

func NewTradingSystemMonitoring(size int) *TradingSystemMonitoring {
	tsa := &TradingSystemMonitoring{}
	tsa.Time          = make([]time.Time, size)
	tsa.GrossProfit   = make([]float64,   size)
	tsa.NetProfit     = make([]float64,   size)
	tsa.GrossDrawdown = make([]float64,   size)
	tsa.NetDrawdown   = make([]float64,   size)

	return tsa
}

//=============================================================================

type PortfolioMonitoringResponse struct {
	BaseMonitoring
	TradingSystems []*TradingSystemMonitoring `json:"tradingSystems"`
}

//=============================================================================
