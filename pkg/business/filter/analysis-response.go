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

import "github.com/bit-fever/portfolio-trader/pkg/core"

//=============================================================================
//===
//=== FilterAnalysisResponse
//===
//=============================================================================

type AnalysisResponse struct {
	TradingSystem TradingSystem   `json:"tradingSystem"`
	Filters       TradingFilters  `json:"filters"`
	Summary       Summary         `json:"summary"`
	Equities      Equities        `json:"equities"`
	Activations   Activations     `json:"activations"`
}

//=============================================================================

type TradingSystem struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

//=============================================================================

type Summary struct {
	UnfProfit       float64 `json:"unfProfit"`
	FilProfit       float64 `json:"filProfit"`
	UnfMaxDrawdown  float64 `json:"unfMaxDrawdown"`
	FilMaxDrawdown  float64 `json:"filMaxDrawdown"`
	UnfWinningPerc  float64 `json:"unfWinningPerc"`
	FilWinningPerc  float64 `json:"filWinningPerc"`
	UnfAverageTrade float64 `json:"unfAverageTrade"`
	FilAverageTrade float64 `json:"filAverageTrade"`
}

//=============================================================================

type Equities struct {
	Days               []int       `json:"days"`
	NetProfits         []float64   `json:"netProfits"`
	UnfilteredEquity   []float64   `json:"unfilteredEquity"`
	FilteredEquity     []float64   `json:"filteredEquity"`
	UnfilteredDrawdown []float64   `json:"unfilteredDrawdown"`
	FilteredDrawdown   []float64   `json:"filteredDrawdown"`
	FilterActivation   []int8      `json:"filterActivation"`
	Average            *core.Plot  `json:"average"`
}

//=============================================================================

type Activation struct {
	Days   []int   `json:"days"`
	Values []int8  `json:"values"`
}

//-----------------------------------------------------------------------------

func (p *Activation) AddDay(day int, value int8) {
	p.Days   = append(p.Days,   day)
	p.Values = append(p.Values, value)
}

//=============================================================================

type Activations struct {
	EquityVsAverage     *Activation `json:"equityVsAverage"`
	PositiveProfit      *Activation `json:"positiveProfit"`
	WinningPercentage   *Activation `json:"winningPercentage"`
	OldVsNew            *Activation `json:"oldVsNew"`
}

//=============================================================================
