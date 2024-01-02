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

package messaging

import (
	"encoding/json"
	"github.com/bit-fever/core/msg"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
	"log/slog"
)

//=============================================================================

func InitMessageListener() {
	slog.Info("Starting message listeners...")

	go msg.ReceiveMessages(msg.QuInventoryUpdatesToPortfolio, handleMessage)
}

//=============================================================================

func handleMessage(m *msg.Message) bool {

	slog.Info("New message received", "origin", m.Origin, "type", m.Type, "source", m.Source)

	if m.Origin == msg.OriginDb {
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
		}
	}

	slog.Error("Dropping message with unknown origin/type!", "origin", m.Origin, "type", m.Type)
	return true
}

//=============================================================================

func setTradingSystem(tsm *TradingSystemMessage, create bool) bool {
	slog.Info("setTradingSystem: Trading system change received", "create", create, "sourceId", tsm.TradingSystem.Id)

	err := db.RunInTransaction(func(tx *gorm.DB) error {
		ts, err := db.GetTradingSystemBySourceId(tx, tsm.TradingSystem.Id)

		if err != nil {
			return err
		}

		if ts == nil {
			ts = &db.TradingSystem{}
			ts.Status = db.TsStatusDisabled
		} else {
			if ts.Username != tsm.TradingSystem.Username {
				slog.Error("Trading system '%v' not owned by user '%v'! Dropping message", tsm.TradingSystem.Id, tsm.TradingSystem.Username)
				return nil
			}
		}

		ts.SourceId        = tsm.TradingSystem.Id
		ts.Username        = tsm.TradingSystem.Username
		ts.WorkspaceCode   = tsm.TradingSystem.WorkspaceCode
		ts.Name            = tsm.TradingSystem.Name
		ts.ProductBrokerId = tsm.TradingSystem.ProductBrokerId
		ts.BrokerSymbol    = tsm.ProductBroker.Symbol
		ts.PointValue      = tsm.ProductBroker.PointValue
		ts.CostPerTrade    = tsm.ProductBroker.CostPerTrade
		ts.MarginValue     = tsm.ProductBroker.MarginValue
		ts.CurrencyId      = tsm.Currency.Id
		ts.CurrencyCode    = tsm.Currency.Code

		if ts.Id == 0 {
			return db.AddTradingSystem(tx, ts)
		}

		db.UpdateTradingSystem(tx, ts)
		return nil
	})

	if err != nil {
		slog.Error("Raised error while processing message")
	} else {
		slog.Info("setTradingSystem: Operation complete")
	}

	return err == nil
}

//=============================================================================
