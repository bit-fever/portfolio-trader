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

	"github.com/tradalia/core/datatype"
)

//=============================================================================
//===
//=== Entities
//===
//=============================================================================

type Portfolio struct {
	Id        uint    `json:"id" gorm:"primaryKey"`
	ParentId  uint    `json:"parentId"`
	Username  string  `json:"username"`
	Name      string  `json:"name"`
}

//=============================================================================

type TsStatus int8

const (
	TsStatusOff     TsStatus = 0
	TsStatusPaused  TsStatus = 1
	TsStatusRunning TsStatus = 2
	TsStatusIdle    TsStatus = 3
	TsStatusBroken  TsStatus = 4
)

//-----------------------------------------------------------------------------

type TsSuggAction int8

const (
	TsActionNone     TsSuggAction = 0
	TsActionTurnOff  TsSuggAction = 1
	TsActionTurnOn   TsSuggAction = 2
	TsActionCheck    TsSuggAction = 3
)

//-----------------------------------------------------------------------------

type TradingSystem struct {
	Id                uint             `json:"id" gorm:"primaryKey"`
	Username          string           `json:"username"`
	Name              string           `json:"name"`
	Timeframe         int              `json:"timeframe"`
	DataProductId     uint             `json:"dataProductId"`
	DataSymbol        string           `json:"dataSymbol"`
	BrokerProductId   uint             `json:"brokerProductId"`
	BrokerSymbol      string           `json:"brokerSymbol"`
	PointValue        float64          `json:"pointValue"`
	CostPerOperation  float64          `json:"costPerOperation"`
	MarginValue       float64          `json:"marginValue"`
	Increment         float64          `json:"increment"`
	MarketType        string           `json:"marketType"`
	CurrencyId        uint             `json:"currencyId"`
	CurrencyCode      string           `json:"currencyCode"`
	CurrencySymbol    string           `json:"currencySymbol"`
	TradingSessionId  uint             `json:"tradingSessionId"`
	SessionName       string           `json:"sessionName"`
	SessionConfig     string           `json:"sessionConfig"`
	AgentProfileId   *uint             `json:"agentProfileId"`
	ExternalRef       string           `json:"externalRef"`
	StrategyType      string           `json:"strategyType"`
	Overnight         bool             `json:"overnight"`
	Tags              string           `json:"tags"`
	Finalized         bool             `json:"finalized"`
	Trading           bool             `json:"trading"`
	Running           bool             `json:"running"`
	AutoActivation    bool             `json:"autoActivation"`
	Active            bool             `json:"active"`
	Status            TsStatus         `json:"status"`
	SuggestedAction   TsSuggAction     `json:"suggestedAction"`
	FirstTrade        *time.Time       `json:"firstTrade"`
	LastTrade         *time.Time       `json:"lastTrade"`
	LastNetProfit     float64          `json:"lastNetProfit"`
	LastNetAvgTrade   float64          `json:"lastNetAvgTrade"`
	LastNumTrades     int              `json:"lastNumTrades"`
	PortfolioId       *uint            `json:"portfolioId"`
	Timezone          string           `json:"timezone"`
	InSampleFrom      datatype.IntDate `json:"inSampleFrom"`
	InSampleTo        datatype.IntDate `json:"inSampleTo"`
	EngineCode        string           `json:"engineCode"`
}

//=============================================================================

type TradingFilter struct {
	TradingSystemId  uint   `json:"tradingSystemId" gorm:"primaryKey"`
	EquAvgEnabled    bool   `json:"equAvgEnabled"`
	EquAvgLen        int    `json:"equAvgLen"`
	PosProEnabled    bool   `json:"posProEnabled"`
	PosProLen        int    `json:"posProLen"`
	WinPerEnabled    bool   `json:"winPerEnabled"`
	WinPerLen        int    `json:"winPerLen"`
	WinPerValue      int    `json:"winPerValue"`
	OldNewEnabled    bool   `json:"oldNewEnabled"`
	OldNewOldLen     int    `json:"oldNewOldLen"`
	OldNewOldPerc    int    `json:"oldNewOldPerc"`
	OldNewNewLen     int    `json:"oldNewNewLen"`
	TrendlineEnabled bool   `json:"trendlineEnabled"`
	TrendlineLen     int    `json:"trendlineLen"`
	TrendlineValue   int    `json:"trendlineValue"`
	DrawdownEnabled  bool   `json:"drawdownEnabled"`
	DrawdownMin      int    `json:"drawdownMin"`
	DrawdownMax      int    `json:"drawdownMax"`
}

//=============================================================================

const (
	TradeTypeLong  = "LO"
	TradeTypeShort = "SH"
	TradeTypeAll   = "**"
)

//-----------------------------------------------------------------------------

type Trade struct {
	Id                 uint       `json:"id" gorm:"primaryKey"`
	TradingSystemId    uint       `json:"tradingSystemId"`
	TradeType          string     `json:"tradeType"`
	EntryDate          *time.Time `json:"entryDate"`
	EntryPrice         float64    `json:"entryPrice"`
	EntryLabel         string     `json:"entryLabel"`
	ExitDate           *time.Time `json:"exitDate"`
	ExitPrice          float64    `json:"exitPrice"`
	ExitLabel          string     `json:"exitLabel"`
	GrossProfit        float64    `json:"grossProfit"`
	Contracts          int        `json:"contracts"`
	EntryDateAtBroker  *time.Time `json:"entryDateAtBroker"`
	EntryPriceAtBroker float64    `json:"entryPriceAtBroker"`
	ExitDateAtBroker   *time.Time `json:"exitDateAtBroker"`
	ExitPriceAtBroker  float64    `json:"exitPriceAtBroker"`
}

//-----------------------------------------------------------------------------

func (t Trade) String() string {
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v|%v|%v",
		t.TradeType,
		t.EntryDate.UTC(), t.EntryPrice, t.EntryLabel,
		t.ExitDate.UTC(),  t.ExitPrice,  t.ExitLabel,
		t.GrossProfit, t.Contracts)
}

//=============================================================================

type DailyReturn struct {
	Id               uint             `json:"id" gorm:"primaryKey"`
	TradingSystemId  uint             `json:"tradingSystemId"`
	Day              datatype.IntDate `json:"day"`
	GrossProfit      float64          `json:"grossProfit"`
	Trades           int              `json:"trades"`
}

//=============================================================================
//===
//=== Table names
//===
//=============================================================================

func (TradingSystem) TableName() string { return "trading_system" }
func (TradingFilter) TableName() string { return "trading_filter" }
func (Trade)         TableName() string { return "trade"          }
func (Portfolio)     TableName() string { return "portfolio"      }
func (DailyReturn)   TableName() string { return "daily_return"   }

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
