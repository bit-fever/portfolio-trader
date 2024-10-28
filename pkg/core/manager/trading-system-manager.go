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

package manager

import (
	"github.com/bit-fever/portfolio-trader/pkg/app"
	"github.com/bit-fever/portfolio-trader/pkg/business/filter"
	"github.com/bit-fever/portfolio-trader/pkg/core"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
	"time"
)

//=============================================================================

func InitTradingSystemManager(cfg *app.Config) *time.Ticker {
	ticker := time.NewTicker(cfg.Update.PeriodHour * time.Hour)

	go func() {
		//--- Wait 2 secs to allow the system to boot properly
		time.Sleep(2 * time.Second)
		run(cfg)

		for range ticker.C {
			run(cfg)
		}
	}()

	return ticker
}

//=============================================================================

func run(cfg *app.Config) {
	slog.Info("Starting trading system manager")
	start := time.Now()

	users, err := getUsersWithTradingSystems()
	if err != nil {
		slog.Error("Cannot get list of users. Update aborted", "error", err)
	} else {
		for _, user := range users {
			err = updateUserSystems(user)
			if err != nil {
				slog.Error("Cannot update trading systems for user", "user", user, "error", err)
			}
		}
	}

	duration := time.Now().Sub(start).Seconds()
	slog.Info("Trading system manager ended", "taskSeconds", duration)
}

//=============================================================================

func getUsersWithTradingSystems() ([]string, error){
	var users []string
	var err error

	err = db.RunInTransaction(func (tx *gorm.DB) error {
		users, err = db.GetUsersWithTradingSystems(tx)
		return err
	})

	return users,err
}

//=============================================================================

func updateUserSystems(user string) error {
	slog.Info("updateUserSystems: Processing trading systems for user", "user", user)

	return db.RunInTransaction(func (tx *gorm.DB) error {
		tss, err := db.GetTradingSystemsByUser(tx, user)
		if err == nil {
			tsMap := createTsMap(tss)
			var allTrades *[]db.Trade
			allTrades,err = db.FindTradesByTsIds(tx, maps.Keys(tsMap))
			if err == nil {
				trMap := createTradeMap(allTrades)
				var filters *[]db.TradingFilter
				filters,err = db.GetTradingFiltersByTsIds(tx, maps.Keys(tsMap))
				if err == nil {
					flMap := createFilterMap(filters)

					for _,ts := range maps.Values(tsMap) {
						list   := trMap[ts.Id]
						filter := flMap[ts.Id]
						err = updateTradingSystem(tx, ts, list, filter)
						if err != nil {
							slog.Info("updateUserSystems: Cannot update trading system","user", user, "tsId", ts.Id, "error", err)
							break
						}
					}
				}

				slog.Info("updateUserSystems: Update ended for user","user", user, "systems", len(*tss), "trades", len(*allTrades))
			}
		}
		return err
	})
}

//=============================================================================

func createTsMap(list *[]db.TradingSystem) map[uint]*db.TradingSystem {
	tsMap := map[uint]*db.TradingSystem{}

	for _, ts := range *list {
		tsMap[ts.Id] = &ts
	}

	return tsMap
}

//=============================================================================

func createTradeMap(list *[]db.Trade) map[uint][]db.Trade {
	trMap := map[uint][]db.Trade{}

	for _, trade := range *list {
		trList := trMap[trade.TradingSystemId]
		trMap[trade.TradingSystemId] = append(trList, trade)
	}

	return trMap
}

//=============================================================================

func createFilterMap(list *[]db.TradingFilter) map[uint]*db.TradingFilter {
	flMap := map[uint]*db.TradingFilter{}

	for _, filter := range *list {
		flMap[filter.TradingSystemId] = &filter
	}

	return flMap
}

//=============================================================================

func updateTradingSystem(tx *gorm.DB, ts *db.TradingSystem, trades []db.Trade, filter *db.TradingFilter) error {
	updateLastMonthStats(ts, trades)
	updateActivationStatus(ts, trades, filter)

	return db.UpdateTradingSystem(tx, ts)
}

//=============================================================================

func updateLastMonthStats(ts *db.TradingSystem, trades []db.Trade) {
	netProit := 0.0
	numTrades:= 0

	startDate := time.Now().Add(-time.Hour * 24*30)

	for _, trade := range trades {
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

func updateActivationStatus(ts *db.TradingSystem, trades []db.Trade, f *db.TradingFilter) {
	if ! ts.Running {
		ts.SuggestedAction = db.TsActionNone
		ts.Status          = db.TsStatusOff
		return
	}

	//--- The trading system is running (i.e. live)

	activValue := false
	if f != nil {
		activValue = filter.CalcActivation(ts, f, trades)
	}

	days := daysToLastTrade(trades)

	switch ts.Activation {
		case db.TsActivationManual:
			handleManualActivation(ts, activValue, days)
		case db.TsActivationAuto:
			handleAutomaticActivation(ts, activValue, days)
		default:
			panic("Unknown filter activation for ts="+ strconv.Itoa(int(ts.Id)))
	}
}

//=============================================================================

func daysToLastTrade(trades []db.Trade) int {
	if len(trades) == 0 {
		//--- If there are no trades yet, we are confident that the first one will arrive soon
		return 0
	}

	lastTrade := trades[len(trades) -1]
	lastDay   := lastTrade.ExitTime

	if lastDay == nil {
		lastDay = lastTrade.EntryTime
	}

	today     := time.Now()

	days := int(today.Sub(*lastDay).Hours() / 24)

	return days
}

//=============================================================================

func handleManualActivation(ts *db.TradingSystem, activValue bool, daysToLastTrade int) {
	if !ts.Active {
		if !activValue {
			ts.Status          = db.TsStatusWaiting
			ts.SuggestedAction = db.TsActionNone
		} else {
			ts.Status          = db.TsStatusWaiting
			ts.SuggestedAction = db.TsActionTurnOn
		}
	} else {
		if !activValue {
			ts.Status          = calcStatus(daysToLastTrade)
			ts.SuggestedAction = db.TsActionTurnOff

			if ts.Status == db.TsStatusBroken {
				ts.SuggestedAction = db.TsActionTurnOffAndCheck
			}
		} else {
			ts.Status          = calcStatus(daysToLastTrade)
			ts.SuggestedAction = db.TsActionNone

			if ts.Status == db.TsStatusBroken {
				ts.SuggestedAction = db.TsActionCheck
			}
		}
	}
}

//=============================================================================

func handleAutomaticActivation(ts *db.TradingSystem, activValue bool, daysToLastTrade int) {
	if !ts.Active {
		if !activValue {
			ts.Status          = db.TsStatusWaiting
			ts.SuggestedAction = db.TsActionNone
		} else {
			ts.Status          = db.TsStatusRunning
			ts.SuggestedAction = db.TsActionNoneTurnedOn
			ts.Active          = true
		}
	} else {
		if !activValue {
			ts.Status          = db.TsStatusWaiting
			ts.SuggestedAction = db.TsActionNoneTurnedOff
			ts.Active          = false
		} else {
			ts.Status          = calcStatus(daysToLastTrade)
			ts.SuggestedAction = db.TsActionNone

			if ts.Status == db.TsStatusBroken {
				ts.SuggestedAction = db.TsActionCheck
			}
		}
	}
}

//=============================================================================

func calcStatus(daysToLastTrade int) db.TsStatus {
	if daysToLastTrade < 7 {
		return db.TsStatusRunning
	}

	if daysToLastTrade < 14 {
		return db.TsStatusIdle
	}

	return db.TsStatusBroken
}

//=============================================================================
