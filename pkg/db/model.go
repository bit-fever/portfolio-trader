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

package db

import "time"

//=============================================================================

type Portfolio struct {
	Id        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

//=============================================================================

type Instrument struct {
	Id        uint   `json:"id" gorm:"primaryKey"`
	Ticker    string `json:"ticker"`
	Name      string `json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

//=============================================================================

type TradingSystem struct {
	Id              uint   `json:"id" gorm:"primaryKey"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	InstrumentId    uint   `json:"instrumentId"`
	PortfolioId     uint   `json:"portfolioId"`
	FirstUpdate     int     `json:"firstUpdate"`
	LastUpdate      int     `json:"lastUpdate"`
	LastPl          float64 `json:"lastPl"`
	TradingDays     int     `json:"tradingDays"`
	NumTrades       int     `json:"numTrades"`
	FilterType      int     `json:"filterType"`
	Filter          string  `json:"filter"`
	SuggestedAction int     `json:"suggestedAction"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//=============================================================================

type TsDailyInfo struct {
	Id              uint    `json:"id" gorm:"primaryKey"`
	TradingSystemId uint    `json:"tradingSystemId"`
	Day             int     `json:"day"`
	OpenProfit      float64 `json:"openProfit"`
	Position        int     `json:"position"`
	GapValue        float64 `json:"gapValue"`
	TrueRange       float64 `json:"trueRange"`
	NumTrades       int     `json:"numTrades"`
}

//=============================================================================

func (Portfolio) TableName() string {
	return "portfolio"
}

//=============================================================================

func (Instrument) TableName() string {
	return "instrument"
}

//=============================================================================

func (TradingSystem) TableName() string {
	return "trading_system"
}

//=============================================================================

func (TsDailyInfo) TableName() string {
	return "ts_daily_info"
}

//=============================================================================
