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

import (
	"github.com/bit-fever/portfolio-trader/pkg/business/filter"
	"github.com/bit-fever/portfolio-trader/pkg/model"
)

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
//=== FilterAnalysisRequest
//===
//=============================================================================

type TradingFilters struct {
	EquAvgEnabled   bool   `json:"equAvgEnabled"`
	EquAvgDays      int    `json:"equAvgDays"`
	PosProEnabled   bool   `json:"posProEnabled"`
	PosProDays      int    `json:"posProDays"`
	WinPerEnabled   bool   `json:"winPerEnabled"`
	WinPerDays      int    `json:"winPerDays"`
	WinPerValue     int    `json:"winPerValue"`
	OldNewEnabled   bool   `json:"oldNewEnabled"`
	OldNewOldDays   int    `json:"oldNewOldDays"`
	OldNewOldPerc   int    `json:"oldNewOldPerc"`
	OldNewNewDays   int    `json:"oldNewNewDays"`
}

//=============================================================================

type FilterAnalysisRequest struct {
	Filters *TradingFilters  `json:"filters,omitempty"`
}

//=============================================================================
//===
//=== FilterAnalysisResponse
//===
//=============================================================================

type TradingSystem struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

//-----------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------

type Equities struct {
	Days               []int       `json:"days"`
	NetProfits         []float64   `json:"netProfits"`
	UnfilteredEquity   []float64   `json:"unfilteredEquity"`
	FilteredEquity     []float64   `json:"filteredEquity"`
	UnfilteredDrawdown []float64   `json:"unfilteredDrawdown"`
	FilteredDrawdown   []float64   `json:"filteredDrawdown"`
	FilterActivation   []int8      `json:"filterActivation"`
	Average            *model.Plot `json:"average"`
}

//-----------------------------------------------------------------------------

type Activations struct {
	EquityVsAverage     *filter.Activation `json:"equityVsAverage"`
	PositiveProfit      *filter.Activation `json:"positiveProfit"`
	WinningPercentage   *filter.Activation `json:"winningPercentage"`
	OldVsNew            *filter.Activation `json:"oldVsNew"`
}

//-----------------------------------------------------------------------------

type FilterAnalysisResponse struct {
	TradingSystem TradingSystem   `json:"tradingSystem"`
	Filters       TradingFilters  `json:"filters"`
	Summary       Summary         `json:"summary"`
	Equities      Equities        `json:"equities"`
	Activations   Activations     `json:"activations"`
}

//=============================================================================
