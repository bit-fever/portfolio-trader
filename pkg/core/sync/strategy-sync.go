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

package sync

import (
	"github.com/bit-fever/core/req"
	"github.com/bit-fever/portfolio-trader/pkg/app"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
	"time"
)

//=============================================================================

func InitPeriodicScan(cfg *app.Config) *time.Ticker {

	ticker := time.NewTicker(cfg.Scan.PeriodHour * time.Hour)

	go func() {
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
	url := cfg.Scan.Address
	slog.Info("Starting to fetch strategies from: "+ url)

	var data []Strategy

	err := req.DoGet(req.GetClient("ws"), url, &data, "")

	if err == nil {
		slog.Info("Got "+ strconv.Itoa(len(data))+ " strategies")

		_ = db.RunInTransaction(func (tx *gorm.DB) error {
			return updateDb(tx, data)
		})
	} else {
		slog.Error("Cannot connect to url. Error is: "+ err.Error())
	}

	slog.Info("Ending to fetch strategies")
}

//=============================================================================

func updateDb(tx *gorm.DB, strategies []Strategy) error {
	slog.Info("Updating database...")

	for _, s := range strategies {
		ts, err := db.GetTradingSystemByName(tx, s.Name)

		if err != nil {
			slog.Error("Cannot scan for trading system '"+ s.Name +"': "+ err.Error())
			return err
		}

		if ts == nil {
			slog.Warn("Trading system '"+ s.Name +"' was not found. Skipping")
			continue
		}

		lastTrade, err := db.FindLastTrade(tx, ts.Id)
		if err != nil {
			slog.Error("Cannot retrieve last trade for trading system '"+ s.Name +"': "+ err.Error())
			return err
		}

		lastTradeTime := time.Unix(0, 0)
		if lastTrade != nil {
			lastTradeTime = *lastTrade.EntryTime
		}

		slog.Info("Updating trading system: "+ s.Name)

		for _, di := range s.DailyInfo {
			tr := createTrade(ts.Id, &di)

			if lastTradeTime.Before(*tr.EntryTime) && tr.GrossProfit != 0 {
				lastTradeTime = *tr.EntryTime
				err = db.AddTrade(tx, tr)
				if err != nil {
					slog.Error("Cannot write trade into database for strategy '"+ s.Name +"': "+ err.Error())
					return err
				}

				//--- Update information on trading system

				if ts.FirstTrade == nil || ts.FirstTrade.After(lastTradeTime) {
					t := lastTradeTime
					ts.FirstTrade = &t
				}

				if ts.LastTrade == nil || ts.LastTrade.Before(lastTradeTime) {
					t := lastTradeTime
					ts.LastTrade = &t
				}
			}
		}

		err = db.UpdateTradingSystem(tx, ts)
		if err != nil {
			return err
		}
	}

	return nil
}

//=============================================================================

func createTrade(tsId uint, di *DailyInfo) *db.Trade {
	y := di.Day / 10000
	m := di.Day / 100 % 100
	d := di.Day % 100

	loc, _ := time.LoadLocation("UTC")

	t := time.Date(y, time.Month(m), d, 8,0,0,0, loc)

	return &db.Trade{
		TradingSystemId: tsId,
		TradeType   : db.TradeTypeLong,
		EntryTime   : &t,
		EntryValue  : 0,
		ExitTime    : &t,
		ExitValue   : 0,
		GrossProfit : di.ClosedProfit,
		NumContracts: 1,
	}
}

//=============================================================================
