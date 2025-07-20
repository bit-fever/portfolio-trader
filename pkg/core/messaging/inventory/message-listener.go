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

package inventory

import (
	"encoding/json"
	"github.com/bit-fever/core/msg"
	"github.com/bit-fever/portfolio-trader/pkg/business"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
	"log/slog"
)

//=============================================================================

func InitMessageListener() {
	slog.Info("Starting inventory message listener...")

	go msg.ReceiveMessages(msg.QuInventoryToPortfolio, handleMessage)
}

//=============================================================================

func handleMessage(m *msg.Message) bool {

	slog.Info("New message received", "source", m.Source, "type", m.Type)

	if m.Source == msg.SourceTradingSystem {
		tsm := TradingSystemMessage{}
		err := json.Unmarshal(m.Entity, &tsm)
		if err != nil {
			slog.Error("Dropping badly formatted message!", "entity", string(m.Entity))
			return true
		}

		if m.Type == msg.TypeCreate {
			return setTradingSystem(&tsm, true)
		}
		if m.Type == msg.TypeUpdate {
			return setTradingSystem(&tsm, false)
		}
		if m.Type == msg.TypeDelete {
			return deleteTradingSystem(&tsm)
		}
	} else if m.Source == msg.SourceDataProduct {
		dpm := DataProductMessage{}
		err := json.Unmarshal(m.Entity, &dpm)
		if err != nil {
			slog.Error("Dropping badly formatted message!", "entity", string(m.Entity))
			return true
		}

		if m.Type == msg.TypeCreate {
			//--- If the broker product is new, there are no trading systems to update. Just return 'true'
			return true
		}

		if m.Type == msg.TypeUpdate {
			return updateDataProduct(&dpm)
		}
	} else if m.Source == msg.SourceBrokerProduct {
		pbm := BrokerProductMessage{}
		err := json.Unmarshal(m.Entity, &pbm)
		if err != nil {
			slog.Error("Dropping badly formatted message!", "entity", string(m.Entity))
			return true
		}

		if m.Type == msg.TypeCreate {
			//--- If the broker product is new, there are no trading systems to update. Just return 'true'
			return true
		}

		if m.Type == msg.TypeUpdate {
			return updateBrokerProduct(&pbm)
		}
	}

	slog.Error("Dropping message with unknown source/type!", "source", m.Source, "type", m.Type)
	return true
}

//=============================================================================

func setTradingSystem(tsm *TradingSystemMessage, create bool) bool {
	slog.Info("setTradingSystem: Trading system change received", "create", create, "id", tsm.TradingSystem.Id)

	err := db.RunInTransaction(func(tx *gorm.DB) error {
		isNew := true

		ts, err := db.GetTradingSystemById(tx, tsm.TradingSystem.Id)

		if err != nil {
			return err
		}

		if ts == nil {
			ts = &db.TradingSystem{}
			ts.Running         = false
			ts.AutoActivation  = false
			ts.Status          = db.TsStatusOff
			ts.Active          = false
			ts.SuggestedAction = db.TsActionNone
		} else {
			isNew = false

			if ts.Username != tsm.TradingSystem.Username {
				slog.Error("Trading system '%v' not owned by user '%v'! Dropping message", tsm.TradingSystem.Id, tsm.TradingSystem.Username)
				return nil
			}
		}

		ts.Id              = tsm.TradingSystem.Id
		ts.Username        = tsm.TradingSystem.Username
		ts.Name            = tsm.TradingSystem.Name
		ts.Timeframe       = tsm.TradingSystem.Timeframe
		ts.DataProductId   = tsm.TradingSystem.DataProductId
		ts.DataSymbol      = tsm.DataProduct.Symbol
		ts.BrokerProductId = tsm.TradingSystem.BrokerProductId
		ts.BrokerSymbol    = tsm.BrokerProduct.Symbol
		ts.PointValue      = tsm.BrokerProduct.PointValue
		ts.CostPerOperation= tsm.BrokerProduct.CostPerOperation
		ts.MarginValue     = tsm.BrokerProduct.MarginValue
		ts.Increment       = tsm.BrokerProduct.Increment
		ts.MarketType      = tsm.BrokerProduct.MarketType
		ts.CurrencyId      = tsm.Currency.Id
		ts.CurrencyCode    = tsm.Currency.Code
		ts.TradingSessionId= tsm.TradingSession.Id
		ts.SessionName     = tsm.TradingSession.Name
		ts.SessionConfig   = tsm.TradingSession.Config
		ts.StrategyType    = tsm.TradingSystem.StrategyType
		ts.Overnight       = tsm.TradingSystem.Overnight
		ts.Tags            = tsm.TradingSystem.Tags
		ts.Finalized       = tsm.TradingSystem.Finalized
		ts.Timezone        = tsm.Exchange.Timezone
		ts.AgentProfileId  = tsm.TradingSystem.AgentProfileId
		ts.ExternalRef     = tsm.TradingSystem.ExternalRef

		err = db.UpdateTradingSystem(tx, ts)

		if err == nil && isNew {
			err = db.SetTradingFilter(tx, &db.TradingFilter{
				TradingSystemId : ts.Id,
			})
		}

		return err
	})

	if err != nil {
		slog.Error("Raised error while processing message")
	} else {
		slog.Info("setTradingSystem: Operation complete")
	}

	return err == nil
}

//=============================================================================

func deleteTradingSystem(tsm *TradingSystemMessage) bool {
	slog.Info("deleteTradingSystem: Trading system deletion received", "id", tsm.TradingSystem.Id)

	err := db.RunInTransaction(func(tx *gorm.DB) error {
		id := tsm.TradingSystem.Id
		return business.DeleteTradingSystem(tx, id)
	})

	if err != nil {
		slog.Error("Raised error while deleting trading system", "error", err.Error())
	} else {
		slog.Info("deleteTradingSystem: Operation complete", "id", tsm.TradingSystem.Id)
	}

	return err == nil
}

//=============================================================================

func updateDataProduct(dpm *DataProductMessage) bool {
	slog.Info("updateDataProduct: Data product change received", "sourceId", dpm.DataProduct.Id)

	err := db.RunInTransaction(func(tx *gorm.DB) error {
		values := map[string]interface{} {
			//--- data_symbol cannot be changed, anyway
			"data_symbol" : dpm.DataProduct.Symbol,
		}

		return db.UpdateDataProductInfo(tx, dpm.DataProduct.Id, values)
	})

	if err != nil {
		slog.Error("Raised error while processing message")
	} else {
		slog.Info("updateDataProduct: Operation complete")
	}

	return err == nil
}

//=============================================================================

func updateBrokerProduct(bpm *BrokerProductMessage) bool {
	slog.Info("updateBrokerProduct: Broker product change received", "sourceId", bpm.BrokerProduct.Id)

	err := db.RunInTransaction(func(tx *gorm.DB) error {
		values := map[string]interface{} {
			"broker_symbol"     : bpm.BrokerProduct.Symbol,
			"point_value"       : bpm.BrokerProduct.PointValue,
			"cost_per_operation": bpm.BrokerProduct.CostPerOperation,
			"margin_value"      : bpm.BrokerProduct.MarginValue,
			"increment"         : bpm.BrokerProduct.Increment,
		}

		return db.UpdateBrokerProductInfo(tx, bpm.BrokerProduct.Id, values)
	})

	if err != nil {
		slog.Error("Raised error while processing message")
	} else {
		slog.Info("updateBrokerProduct: Operation complete")
	}

	return err == nil
}

//=============================================================================
