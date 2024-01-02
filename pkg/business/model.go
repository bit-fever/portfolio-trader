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

//=============================================================================
//===
//=== Portfolio monitoring
//===
//=============================================================================

type PortfolioMonitoringParams struct {
	TsIds  []uint `form:"tsIds"  binding:"required,min=1,dive"`
	Period    int `form:"period" binding:"required,min=1,max=5000"`
}

//=============================================================================

type BaseMonitoring struct {
	Days        []int      `json:"days"`
	RawProfit   []float64  `json:"rawProfit"`
	NetProfit   []float64  `json:"netProfit"`
	RawDrawdown []float64  `json:"rawDrawdown"`
	NetDrawdown []float64  `json:"netDrawdown"`
	NumTrades   []int      `json:"numTrades"`
}

//-----------------------------------------------------------------------------

type TradingSystemMonitoring struct {
	BaseMonitoring
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

//-----------------------------------------------------------------------------

type PortfolioMonitoringResponse struct {
	BaseMonitoring
	TradingSystems []*TradingSystemMonitoring `json:"tradingSystems"`
}

//-----------------------------------------------------------------------------

func NewTradingSystemAnalysis(size int) *TradingSystemMonitoring {
	tsa := &TradingSystemMonitoring{}
	tsa.Days        = make([]int,     size)
	tsa.RawProfit   = make([]float64, size)
	tsa.NetProfit   = make([]float64, size)
	tsa.RawDrawdown = make([]float64, size)
	tsa.NetDrawdown = make([]float64, size)
	tsa.NumTrades   = make([]int,     size)

	return tsa
}

//=============================================================================
//===
//=== Filtering analysis
//===
//=============================================================================

type LongShort struct {
	Enabled      bool    `json:"enabled"`
	LongPeriod   int     `json:"longPeriod"`
	ShortPeriod  int     `json:"shortPeriod"`
	Threshold    float32 `json:"threshold"`
	ShortPosPerc int     `json:"shortPosPerc"`
}

//-----------------------------------------------------------------------------

type EquityAverage struct {
	Enabled bool `json:"enabled"`
	Days    int  `json:"days"`
}

//-----------------------------------------------------------------------------

type FilteringConfig struct {
	LongShort     `json:"longShort"`
	EquityAverage `json:"equityAverage"`
}

//-----------------------------------------------------------------------------

type FilteringParams struct {
	NoConfig bool   `json:"noConfig"`
	FilteringConfig `json:"config"`
}

//=============================================================================

type TradingSystem struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

//-----------------------------------------------------------------------------

type Summary struct {
	UnfilteredProfit float64 `json:"unfilteredProfit"`
	FilteredProfit   float64 `json:"filteredProfit"`
	UnfilteredMaxDD  float64 `json:"unfilteredMaxDD"`
	FilteredMaxDD    float64 `json:"filteredMaxDD"`
}

//-----------------------------------------------------------------------------

type Equities struct {
	Days               []int      `json:"days"`
	UnfilteredProfit   []float64  `json:"unfilteredProfit"`
	FilteredProfit     []float64  `json:"filteredProfit"`
	UnfilteredDrawdown []float64  `json:"unfilteredDrawdown"`
	FilteredDrawdown   []float64  `json:"filteredDrawdown"`
	Average            []float64  `json:"average"`
}

//-----------------------------------------------------------------------------

type FilteringResponse struct {
	TradingSystem      `json:"tradingSystem"`
	Summary            `json:"summary"`
	Equities           `json:"equities"`
	FilteringConfig    `json:"config"`
}

//=============================================================================
