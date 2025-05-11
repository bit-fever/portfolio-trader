//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

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

package performance

import (
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"time"
)

//=============================================================================

type Value struct {
	Total float64 `json:"total"`
	Long  float64 `json:"long"`
	Short float64 `json:"short"`
}

//=============================================================================

type Performance struct {
	Profit        Value `json:"profit"`
	MaxDrawdown   Value `json:"maxDrawdown"`
	AverageTrade  Value `json:"averageTrade"`
	PercentProfit Value `json:"percentProfit"`
}

//=============================================================================

type Equities struct {
	Time          *[]time.Time `json:"time"`
	GrossEquity   *[]float64   `json:"grossProfit"`
	NetEquity     *[]float64   `json:"netProfit"`
	GrossDrawdown *[]float64   `json:"grossDrawdown"`
	NetDrawdown   *[]float64   `json:"netDrawdown"`
	Trades        int          `json:"trades"`
}

//=============================================================================

type AnalysisResponse struct {
	TradingSystem     *db.TradingSystem `json:"tradingSystem"`
	Gross             Performance       `json:"gross"`
	Net               Performance       `json:"net"`
	AllEquities       *Equities         `json:"allEquities"`
	LongEquities      *Equities         `json:"longEquities"`
	ShortEquities     *Equities         `json:"shortEquities"`
	Trades            *[]db.Trade       `json:"trades"`
	//Periodic          int
	//AvgBarsInTrade    int               `json:"avgBarsInTrade"`
	//AvgBarsBetwTrades int               `json:"avgBarsBetwTrades"`
}

//=============================================================================
