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

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

//=============================================================================
//===
//=== Entities
//===
//=============================================================================

const TsStatusEnabled  = "en"
const TsStatusDisabled = "di"

//=============================================================================

type TradingSystem struct {
	Id               uint        `json:"id" gorm:"primaryKey"`
	Username         string      `json:"username"`
	WorkspaceCode    string      `json:"workspaceCode"`
	Name             string      `json:"name"`
	Status           string      `json:"status"`
	FirstTrade       *time.Time  `json:"firstTrade"`
	LastTrade        *time.Time  `json:"lastTrade"`
	LmNetProfit      float64     `json:"lmNetProfit"`
	LmAvgTrade       float64     `json:"lmAvgTrade"`
	LmNumTrades      int         `json:"lmNumTrades"`
	BrokerProductId  uint        `json:"brokerProductId"`
	BrokerSymbol     string      `json:"brokerSymbol"`
	PointValue       float32     `json:"pointValue"`
	CostPerOperation float32     `json:"costPerOperation"`
	MarginValue      float32     `json:"marginValue"`
	Increment        float64     `json:"increment"`
	CurrencyId       uint        `json:"currencyId"`
	CurrencyCode     string      `json:"currencyCode"`
}

//=============================================================================

type TradingFilter struct {
	TradingSystemId uint   `json:"tradingSystemId" gorm:"primaryKey"`
	EquAvgEnabled   bool   `json:"equAvgEnabled"`
	EquAvgLen       int    `json:"equAvgLen"`
	PosProEnabled   bool   `json:"posProEnabled"`
	PosProLen       int    `json:"posProLen"`
	WinPerEnabled   bool   `json:"winPerEnabled"`
	WinPerLen       int    `json:"winPerLen"`
	WinPerValue     int    `json:"winPerValue"`
	OldNewEnabled   bool   `json:"shoLonEnabled"`
	OldNewOldLen    int    `json:"oldNewOldLen"`
	OldNewOldPerc   int    `json:"oldNewOldPerc"`
	OldNewNewLen    int    `json:"oldNewNewLen"`
}

//=============================================================================

type TradeType int8

const (
	TradeTypeLong  TradeType = 0
	TradeTypeShort TradeType = 1
)

type Trade struct {
	Id              uint       `json:"id" gorm:"primaryKey"`
	TradingSystemId uint       `json:"tradingSystemId"`
	TradeType       TradeType  `json:"tradeType"`
	EntryTime       *time.Time `json:"entryTime"`
	EntryValue      float64    `json:"entryValue"`
	ExitTime        *time.Time `json:"exitTime"`
	ExitValue       float64    `json:"exitValue"`
	GrossProfit     float64    `json:"grossProfit"`
	NumContracts    int        `json:"numContracts"`
}

//=============================================================================
//===
//=== Table names
//===
//=============================================================================

func (TradingSystem) TableName() string { return "trading_system" }
func (TradingFilter) TableName() string { return "trading_filter" }
func (Trade)         TableName() string { return "trade"          }

//=============================================================================
//===
//=== ParamMap type
//===
//=============================================================================

type ParamMap map[string]any

//=============================================================================

func (pm *ParamMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	err := json.Unmarshal(bytes, pm)

	return err
}

//=============================================================================

func (pm ParamMap) Value() (driver.Value, error) {
	return json.Marshal(pm)
}

//=============================================================================
