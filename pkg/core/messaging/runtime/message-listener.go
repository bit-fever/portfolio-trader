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
	"log/slog"
	"sort"
	"time"

	"github.com/bit-fever/core/datatype"
	"github.com/bit-fever/core/msg"
	"github.com/bit-fever/portfolio-trader/pkg/business/filter"
	"github.com/bit-fever/portfolio-trader/pkg/consts"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
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
		tm := TradeListMessage{}
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

func handleNewTrades(tm *TradeListMessage) bool {
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

		if !ts.Trading {
			slog.Warn("handleNewTrades: Trading system is not in TRADING mode. Discarding trades", "id", tsId)
			return nil
		}

		var trades *[]db.Trade
		trades,err = db.FindTradesByTradingSystemId(tx, tsId)
		if err == nil {
			var dailyProfits *[]db.DailyReturn
			dailyProfits,err = db.FindDailyReturnsByTradingSystemId(tx, tsId)
			if err == nil {
				var tf *db.TradingFilter
				tf, err = db.GetTradingFilterByTsId(tx, tsId)
				if err == nil {
					trades,err = addNewTrades(tx, ts, trades, tm.Trades)
					if err == nil {
						err = addNewDailyProfits(tx, ts, dailyProfits, tm.DailyProfits)
						if err == nil {
							err = updateTradingSystem(tx, ts, trades, tf)
						}
					}
				}
			}
		}

		return err
	})

	return err == nil
}

//=============================================================================

func addNewTrades(tx *gorm.DB, ts *db.TradingSystem, trades *[]db.Trade, newTrades []*TradeItem) (*[]db.Trade, error) {
	list := *trades

	tradeSet := map[string]bool{}
	for _, dbt := range *trades {
		tradeSet[dbt.String()] = true
	}

	firstTrade := ts.FirstTrade
	lastTrade  := ts.LastTrade

	for _, tr := range newTrades {
		dbTr := toDbTrade(ts.Id, tr)
		_, exists := tradeSet[dbTr.String()]
		if exists {
			continue
		}

		//--- We need to add trades that are outside of [firstTrade .. lastTrade]
		//--- because we will have duplicates when importing from external strategies
		//--- Example: we have @NQ and we run the strategy on the full period to get lots of data.
		//--- Then, when switching to live, the instrument will switch to something like @NQM25 for roughly
		//--- 180 days. @NQ and @NQM25 are slightly different and there will be a jump in the continuous
		//--- contract causing duplicates between @NQ and @NQM25 during the last 180 days

		if isTheTradeOutOfRange(ts, dbTr) {
			tradeSet[dbTr.String()] = true
			err  := db.AddTrade(tx, dbTr)

			if err != nil {
				return nil, err
			}

			list = append(list, *dbTr)

			//--- Update information on trading system
			//--- It is better to use the exit date for first/last trade because a trade could last
			//--- for 7+ days and the IDLE flag is impacted

			if firstTrade == nil || firstTrade.After(*tr.ExitDate) {
				firstTrade = tr.ExitDate
			}

			if lastTrade == nil || lastTrade.Before(*tr.ExitDate) {
				lastTrade = tr.ExitDate
			}
		}
	}

	ts.FirstTrade = firstTrade
	ts.LastTrade  = lastTrade

	//--- Sort final list as new trades could be in the past

	sort.Slice(list, func(i,j int) bool {
		return list[i].EntryDate.Before(*list[j].EntryDate)
	})

	return &list, nil
}

//=============================================================================

func isTheTradeOutOfRange(ts *db.TradingSystem, t *db.Trade) bool {
	if ts.FirstTrade == nil || ts.LastTrade == nil {
		return true
	}

	//--- Better to use the exit date (please, see above)
	return t.ExitDate.Before(*ts.FirstTrade) || t.ExitDate.After(*ts.LastTrade)
}

//=============================================================================

func toDbTrade(tsId uint, t *TradeItem) *db.Trade {
	return &db.Trade{
		TradingSystemId: tsId,
		TradeType      : t.TradeType,
		EntryDate      : t.EntryDate,
		EntryPrice     : t.EntryPrice,
		EntryLabel     : t.EntryLabel,
		ExitDate       : t.ExitDate,
		ExitPrice      : t.ExitPrice,
		ExitLabel      : t.ExitLabel,
		GrossProfit    : t.GrossProfit,
		Contracts      : t.Contracts,
	}
}

//=============================================================================

func addNewDailyProfits(tx *gorm.DB, ts *db.TradingSystem, profits *[]db.DailyReturn, newProfits []*DailyProfitItem) error {
	profitSet := map[datatype.IntDate]bool{}
	for _, dp := range *profits {
		profitSet[dp.Day] = true
	}

	for _, dp := range newProfits {
		dbDp := toDbDailyProfit(ts.Id, dp)
		_, exists := profitSet[dbDp.Day]
		if !exists {
			profitSet[dbDp.Day] = true
			err  := db.AddDailyReturn(tx, dbDp)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

//=============================================================================

func toDbDailyProfit(tsId uint, p *DailyProfitItem) *db.DailyReturn {
	return &db.DailyReturn{
		TradingSystemId: tsId,
		Day            : p.Day,
		GrossProfit    : p.GrossProfit,
		Trades         : p.Trades,
	}
}

//=============================================================================

func updateTradingSystem(tx *gorm.DB, ts *db.TradingSystem, trades *[]db.Trade, filter *db.TradingFilter) error {
	updateActivationStatus(ts, trades, filter)

	//--- If we got new trades, probably we have to set an idle/broken state to running

	if ts.Status == db.TsStatusIdle || ts.Status == db.TsStatusBroken {
		idleStart := time.Now().Add(-time.Hour * 24 * time.Duration(consts.IdleDays))

		if ts.LastTrade.After(idleStart) {
			ts.Status = db.TsStatusRunning
		}
	}

	return db.UpdateTradingSystem(tx, ts)
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

	if ts.AutoActivation {
		handleAutomaticActivation(ts, activValue)
	} else {
		handleManualActivation(ts, activValue)
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
		} else {
			ts.SuggestedAction = db.TsActionNone
		}
	}
}

//=============================================================================

func handleAutomaticActivation(ts *db.TradingSystem, activValue bool) {
	ts.SuggestedAction = db.TsActionNone

	if !ts.Active {
		if activValue {
			ts.Status = db.TsStatusRunning
			ts.Active = true
			activate(ts)
			notifyRuntime(ts)
		}
	} else {
		if !activValue {
			ts.Status = db.TsStatusPaused
			ts.Active = false
			activate(ts)
			notifyRuntime(ts)
		}
	}
}

//=============================================================================

func activate(ts *db.TradingSystem) {
	//TODO
}

//=============================================================================

func notifyRuntime(ts *db.TradingSystem) {
	//TODO
}

//=============================================================================
