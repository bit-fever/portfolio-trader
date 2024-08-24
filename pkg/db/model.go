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
	Id              uint       `json:"id" gorm:"primaryKey"`
	Username        string     `json:"username"`
	WorkspaceCode   string     `json:"workspaceCode"`
	Name            string     `json:"name"`
	Status          string     `json:"status"`
	FirstUpdate     int        `json:"firstUpdate"`
	LastUpdate      int        `json:"lastUpdate"`
	ClosedProfit    float64    `json:"closedProfit"`
	TradingDays     int        `json:"tradingDays"`
	NumTrades       int        `json:"numTrades"`
	BrokerProductId uint       `json:"brokerProductId"`
	BrokerSymbol    string     `json:"brokerSymbol"`
	PointValue      float32    `json:"pointValue"`
	CostPerTrade    float32    `json:"costPerTrade"`
	MarginValue     float32    `json:"marginValue"`
	CurrencyId      uint       `json:"currencyId"`
	CurrencyCode    string     `json:"currencyCode"`
}

//=============================================================================

type TradingFilters struct {
	TradingSystemId uint   `json:"tradingSystemId" gorm:"primaryKey"`
	EquAvgEnabled   bool   `json:"equAvgEnabled"`
	EquAvgDays      int    `json:"equAvgDays"`
	PosProEnabled   bool   `json:"posProEnabled"`
	PosProDays      int    `json:"posProDays"`
	WinPerEnabled   bool   `json:"winPerEnabled"`
	WinPerDays      int    `json:"winPerDays"`
	WinPerValue     int    `json:"winPerValue"`
	OldNewEnabled   bool   `json:"shoLonEnabled"`
	OldNewOldDays   int    `json:"oldNewOldDays"`
	OldNewOldPerc   int    `json:"oldNewOldPerc"`
	OldNewNewDays   int    `json:"oldNewNewDays"`
}

//=============================================================================

type DailyInfo struct {
	Id              uint    `json:"id" gorm:"primaryKey"`
	TradingSystemId uint    `json:"tradingSystemId"`
	Day             int     `json:"day"`
	OpenProfit      float64 `json:"openProfit"`
	ClosedProfit    float64 `json:"closedProfit"`
	Position        int     `json:"position"`
	NumTrades       int     `json:"numTrades"`
}

//=============================================================================
//===
//=== Table names
//===
//=============================================================================

func (TradingSystem)  TableName() string { return "trading_system" }
func (TradingFilters) TableName() string { return "trading_filter" }
func (DailyInfo)      TableName() string { return "daily_info"     }

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
