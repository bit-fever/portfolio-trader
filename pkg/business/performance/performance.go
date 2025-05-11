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
	"github.com/bit-fever/portfolio-trader/pkg/core"
	"github.com/bit-fever/portfolio-trader/pkg/db"
)

//=============================================================================

func GetPerformanceAnalysis(ts *db.TradingSystem, trades *[]db.Trade) *AnalysisResponse {
	res := AnalysisResponse{}
	res.TradingSystem = ts
	res.Trades        = trades

	allEq  , allMaxGrossDD  , allMaxNetDD   := calcEquities(ts, trades, db.TradeTypeAll)
	longEq , longMaxGrossDD , longMaxNetDD  := calcEquities(ts, trades, db.TradeTypeLong)
	shortEq, shortMaxGrossDD, shortMaxNetDD := calcEquities(ts, trades, db.TradeTypeShort)

	res.AllEquities   = allEq
	res.LongEquities  = longEq
	res.ShortEquities = shortEq

	res.Gross.MaxDrawdown.Total = allMaxGrossDD
	res.Gross.MaxDrawdown.Long  = longMaxGrossDD
	res.Gross.MaxDrawdown.Short = shortMaxGrossDD
	res.Net  .MaxDrawdown.Total = allMaxNetDD
	res.Net  .MaxDrawdown.Long  = longMaxNetDD
	res.Net  .MaxDrawdown.Short = shortMaxNetDD

	res.Gross.Profit.Total = calcProfit(allEq  .GrossEquity)
	res.Gross.Profit.Long  = calcProfit(longEq .GrossEquity)
	res.Gross.Profit.Short = calcProfit(shortEq.GrossEquity)
	res.Net  .Profit.Total = calcProfit(allEq  .NetEquity)
	res.Net  .Profit.Long  = calcProfit(longEq .NetEquity)
	res.Net  .Profit.Short = calcProfit(shortEq.NetEquity)

	res.Gross.AverageTrade.Total = calcAvgTrade(res.Gross.Profit.Total, allEq.Trades)
	res.Gross.AverageTrade.Long  = calcAvgTrade(res.Gross.Profit.Long , longEq.Trades)
	res.Gross.AverageTrade.Short = calcAvgTrade(res.Gross.Profit.Short, shortEq.Trades)
	res.Net  .AverageTrade.Total = calcAvgTrade(res.Net  .Profit.Total, allEq.Trades)
	res.Net  .AverageTrade.Long  = calcAvgTrade(res.Net  .Profit.Long , longEq.Trades)
	res.Net  .AverageTrade.Short = calcAvgTrade(res.Net  .Profit.Short, shortEq.Trades)

	return &res
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func calcEquities(ts *db.TradingSystem, trades *[]db.Trade, tradeType string) (*Equities, float64, float64) {
	timeSlice, grossProfits := core.BuildGrossProfits(trades, tradeType)
	netProfits              := core.BuildNetProfits(grossProfits, ts.CostPerOperation)

	grossEquity := core.BuildEquity(grossProfits)
	netEquity   := core.BuildEquity(netProfits)

	grossDD, maxGrossDD := core.BuildDrawDown(grossEquity)
	netDD  , maxNetDD   := core.BuildDrawDown(netEquity)

	return &Equities{
		Time         : timeSlice,
		GrossEquity  : grossEquity,
		NetEquity    : netEquity,
		GrossDrawdown: grossDD,
		NetDrawdown  : netDD,
		Trades       : len(*timeSlice),
	}, maxGrossDD, maxNetDD
}

//=============================================================================

func calcProfit(equity *[]float64) float64 {
	if equity == nil || len(*equity)==0 {
		return 0
	}

	return (*equity)[len(*equity) -1]
}

//=============================================================================

func calcAvgTrade(value float64, count int) float64 {
	if count == 0 {
		return 0
	}

	return core.Trunc2d(value / float64(count))
}

//=============================================================================
