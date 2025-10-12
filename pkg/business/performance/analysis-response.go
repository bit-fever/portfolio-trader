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
	"time"

	"github.com/bit-fever/core/datatype"
	"github.com/bit-fever/portfolio-trader/pkg/core"
	"github.com/bit-fever/portfolio-trader/pkg/core/stats"
	"github.com/bit-fever/portfolio-trader/pkg/db"
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

type General struct {
	FromDate datatype.IntDate  `json:"fromDate"`
	ToDate   datatype.IntDate  `json:"toDate"`
}

//=============================================================================

type Aggregates struct {
	Annual *[]*AnnualAggregate `json:"annual"`
}

//=============================================================================

type AnnualAggregate struct {
	Year          int     `json:"year"`
	GrossProfit   float64 `json:"grossProfit"`
	GrossAvgTrade float64 `json:"grossAvgTrade"`
	GrossWinPerc  float64 `json:"grossWinPerc"`
	NetProfit     float64 `json:"netProfit"`
	NetAvgTrade   float64 `json:"netAvgTrade"`
	NetWinPerc    float64 `json:"netWinPerc"`
	Trades        int     `json:"trades"`
}

//-----------------------------------------------------------------------------

func NewAggregate(tr *db.Trade, cost float64) *AnnualAggregate{
	a := &AnnualAggregate{
		Year         : tr.ExitDate.Year(),
		GrossProfit  : tr.GrossProfit,
		GrossAvgTrade: 0,
		GrossWinPerc : 0,
		NetProfit    : tr.GrossProfit - 2 * cost,
		NetAvgTrade  : 0,
		NetWinPerc   : 0,
		Trades       : 1,
	}

	if a.GrossProfit > 0 {
		a.GrossWinPerc = 1
	}

	if a.NetProfit > 0 {
		a.NetWinPerc = 1
	}

	return a
}

//-----------------------------------------------------------------------------

func (a *AnnualAggregate) addTrade(tr *db.Trade, cost float64) {
	netProfit := tr.GrossProfit - 2 * cost

	a.GrossProfit += tr.GrossProfit
	a.NetProfit   += netProfit
	a.Trades++

	if tr.GrossProfit > 0 {
		a.GrossWinPerc++
	}

	if netProfit > 0 {
		a.NetWinPerc++
	}
}

//-----------------------------------------------------------------------------

func (a *AnnualAggregate) consolidate() {
	a.GrossAvgTrade = core.Trunc2d(a.GrossProfit  / float64(a.Trades))
	a.GrossWinPerc  = core.Trunc2d(a.GrossWinPerc / float64(a.Trades) * 100)
	a.NetAvgTrade   = core.Trunc2d(a.NetProfit    / float64(a.Trades))
	a.NetWinPerc    = core.Trunc2d(a.NetWinPerc   / float64(a.Trades) * 100)
}

//=============================================================================

type TradeDistribution struct {
	SharpeRatioAnnualized Value `json:"sharpeRatioAnnualized"`
	StandardDevAnnualized Value `json:"standardDevAnnualized"`
	LowerTail             Value `json:"lowerTail"`
	UpperTail             Value `json:"upperTail"`
}

//=============================================================================

type Distribution struct {
	Mean           float64           `json:"mean"`
	Median         float64           `json:"median"`
	StandardDev    float64           `json:"standardDev"`
	SharpeRatio    float64           `json:"sharpeRatio"`
	LowerTail      float64           `json:"lowerTail"`
	UpperTail      float64           `json:"upperTail"`
	Skewness       float64           `json:"skewness"`
	Histogram      *stats.Histogram  `json:"histogram"`
}

//=============================================================================

type Distributions struct {
	Daily             *Distribution `json:"daily"`
	TradesAllGross    *Distribution `json:"tradesAllGross"`
	TradesAllNet      *Distribution `json:"tradesAllNet"`
	TradesLongGross   *Distribution `json:"tradesLongGross"`
	TradesLongNet     *Distribution `json:"tradesLongNet"`
	TradesShortGross  *Distribution `json:"tradesShortGross"`
	TradesShortNet    *Distribution `json:"tradesShortNet"`
	AnnualSharpeRatio float64       `json:"annualSharpeRatio"`
	AnnualStandardDev float64       `json:"annualStandardDev"`
}

//=============================================================================

type AnalysisResponse struct {
	General           General           `json:"general"`
	TradingSystem     *db.TradingSystem `json:"tradingSystem"`
	Gross             Performance       `json:"gross"`
	Net               Performance       `json:"net"`
	AllEquities       *Equities         `json:"allEquities"`
	LongEquities      *Equities         `json:"longEquities"`
	ShortEquities     *Equities         `json:"shortEquities"`
	Trades            *[]db.Trade       `json:"trades"`
	Aggregates        Aggregates        `json:"aggregates"`
	Distributions     Distributions     `json:"distributions"`
}

//=============================================================================
