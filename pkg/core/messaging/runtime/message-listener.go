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
	"encoding/json"
	"github.com/bit-fever/core/msg"
	"github.com/bit-fever/portfolio-trader/pkg/business/filter"
	"github.com/bit-fever/portfolio-trader/pkg/core"
	"github.com/bit-fever/portfolio-trader/pkg/core/tradingsystem"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
	"time"
)

//=============================================================================

func InitMessageListener() {
	slog.Info("Starting runtime message listener...")

	go msg.ReceiveMessages(msg.QuRuntimeToPortfolio, handleMessage)
}

//=============================================================================

func handleMessage(m *msg.Message) bool {

	slog.Info("New message received", "source", m.Source, "type", m.Type)

	if m.Source == msg.SourceTrade {
		tm := TradeMessage{}
		err := json.Unmarshal(m.Entity, &tm)
		if err != nil {
			slog.Error("Dropping badly formatted message!", "entity", string(m.Entity))
			return true
		}

		if m.Type == msg.TypeCreate {
			return handleNewTrades(&tm)
		}
	}

	slog.Error("Dropping message with unknown source/type!", "source", m.Source, "type", m.Type)
	return true
}

//=============================================================================

func handleNewTrades(tm *TradeMessage) bool {
	tsId := tm.TradingSystemId

	slog.Info("handleNewTrades: Processing new trades for trading systems", "id", tsId)

	err := db.RunInTransaction(func (tx *gorm.DB) error {
		ts, err := db.GetTradingSystemById(tx, tsId)
		if err != nil {
			slog.Error("handleNewTrades: Cannot retrieve trading system", "id", tsId, "error", err.Error())
			return err
		}

		if ts == nil {
			slog.Error("handleNewTrades: Trading system not found. Discarding trades", "id", tsId)
			return nil
		}

		var trades *[]db.Trade
		trades,err = db.FindTradesByTsId(tx, tsId)
		if err == nil {
			var tf *db.TradingFilter
			tf, err = db.GetTradingFilterByTsId(tx, tsId)
			if err == nil {
				trades,err = addNewTrades(tx, ts, trades, tm.Trades)
				if err == nil {
					err = updateTradingSystem(tx, ts, trades, tf)
				}
			}
		}

		return err
	})

	return err == nil
}

//=============================================================================

func addNewTrades(tx *gorm.DB, ts *db.TradingSystem, trades *[]db.Trade, newTrades []*Trade) (*[]db.Trade, error) {
	list := *trades

	for _, tr := range newTrades {
		dbTr := toDbTrade(ts.Id, tr)
		err  := db.AddTrade(tx, dbTr)

		if err != nil {
			return nil, err
		}

		list = append(list, *dbTr)

		//--- Update information on trading system

		if ts.FirstTrade == nil || ts.FirstTrade.After(*tr.EntryTime) {
			ts.FirstTrade = tr.EntryTime
		}

		if ts.LastTrade == nil || ts.LastTrade.Before(*tr.EntryTime) {
			ts.LastTrade = tr.EntryTime
		}
	}

	return &list, nil
}

//=============================================================================

func toDbTrade(tsId uint, t *Trade) *db.Trade {
	return &db.Trade{
		TradingSystemId: tsId,
		TradeType      : t.TradeType,
		EntryTime      : t.EntryTime,
		EntryValue     : t.EntryValue,
		ExitTime       : t.ExitTime,
		ExitValue      : t.ExitValue,
		GrossProfit    : t.GrossProfit,
		NumContracts   : t.NumContracts,
	}
}

//=============================================================================

func updateTradingSystem(tx *gorm.DB, ts *db.TradingSystem, trades *[]db.Trade, filter *db.TradingFilter) error {
	updateLastMonthStats(ts, trades)
	updateActivationStatus(ts, trades, filter)

	//--- If we got new trades, probably we have to set an idle/broken state to running

	if ts.Status == db.TsStatusIdle || ts.Status == db.TsStatusBroken {
		idleStart := time.Now().Add(-time.Hour * 24 * time.Duration(tradingsystem.IdleDays))

		if ts.LastTrade.After(idleStart) {
			ts.Status = db.TsStatusRunning
		}
	}

	now := time.Now()
	ts.LastUpdate = &now

	return db.UpdateTradingSystem(tx, ts)
}

//=============================================================================

func updateLastMonthStats(ts *db.TradingSystem, trades *[]db.Trade) {
	netProit := 0.0
	numTrades:= 0

	startDate := time.Now().Add(-time.Hour * 24*30)

	for _, trade := range *trades {
		if trade.ExitTime.After(startDate) {
			netProit += trade.GrossProfit - 2 * float64(ts.CostPerOperation)
			numTrades++
		}
	}

	ts.LmNetProfit   = netProit
	ts.LmNumTrades   = numTrades
	ts.LmNetAvgTrade = 0

	if numTrades != 0 {
		ts.LmNetAvgTrade = core.Trunc2d(netProit / float64(numTrades))
	}
}

//=============================================================================

func updateActivationStatus(ts *db.TradingSystem, trades *[]db.Trade, f *db.TradingFilter) {
	if ! ts.Running {
		ts.SuggestedAction = db.TsActionNone
		ts.Status          = db.TsStatusOff
		return
	}

	//--- The trading system is running (i.e. live)

	activValue := false
	if f != nil {
		activValue = filter.CalcActivation(ts, f, *trades)
	}

	switch ts.Activation {
		case db.TsActivationManual:
			handleManualActivation(ts, activValue)
		case db.TsActivationAuto:
			handleAutomaticActivation(ts, activValue)
		default:
			panic("Unknown filter activation for ts="+ strconv.Itoa(int(ts.Id)))
	}
}

//=============================================================================

func handleManualActivation(ts *db.TradingSystem, activValue bool) {
	if !ts.Active {
		if !activValue {
			ts.SuggestedAction = db.TsActionNone
		} else {
			ts.SuggestedAction = db.TsActionTurnOn
		}
	} else {
		if !activValue {
			ts.SuggestedAction = db.TsActionTurnOff

			if ts.Status == db.TsStatusBroken {
				ts.SuggestedAction = db.TsActionTurnOffAndCheck
			}
		} else {
			ts.SuggestedAction = db.TsActionNone

			if ts.Status == db.TsStatusBroken {
				ts.SuggestedAction = db.TsActionCheck
			}
		}
	}
}

//=============================================================================

func handleAutomaticActivation(ts *db.TradingSystem, activValue bool) {
	if !ts.Active {
		if !activValue {
			ts.SuggestedAction = db.TsActionNone
		} else {
			ts.Status          = db.TsStatusRunning
			ts.SuggestedAction = db.TsActionNoneTurnedOn
			ts.Active          = true
			notifyRuntime(ts)
		}
	} else {
		if !activValue {
			ts.Status          = db.TsStatusPaused
			ts.SuggestedAction = db.TsActionNoneTurnedOff
			ts.Active          = false
			notifyRuntime(ts)
		} else {
			ts.SuggestedAction = db.TsActionNone

			if ts.Status == db.TsStatusBroken {
				ts.SuggestedAction = db.TsActionCheck
			}
		}
	}
}

//=============================================================================

func notifyRuntime(ts *db.TradingSystem) {

}

//=============================================================================
