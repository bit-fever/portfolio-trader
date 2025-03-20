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

package runtime

import (
	"time"
)

//=============================================================================

type TradeListMessage struct {
	TradingSystemId uint         `json:"tradingSystemId"`
	Trades          []*TradeItem `json:"trades"`
}

//=============================================================================

type TradeItem struct {
	TradeType       string       `json:"tradeType"`
	EntryDate       *time.Time   `json:"entryDate"`
	EntryPrice      float64      `json:"entryPrice"`
	EntryLabel      string       `json:"entryLabel"`
	ExitDate        *time.Time   `json:"exitDate"`
	ExitPrice       float64      `json:"exitPrice"`
	ExitLabel       string       `json:"exitLabel"`
	GrossProfit     float64      `json:"grossProfit"`
	Contracts       int          `json:"contracts"`
}

//=============================================================================
