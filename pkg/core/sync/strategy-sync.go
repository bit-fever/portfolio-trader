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

	var data []strategy

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

func updateDb(tx *gorm.DB, inStrategies []strategy) error {
	slog.Info("Updating database...")

	for _, s := range inStrategies {
		ts, err := db.GetTradingSystemByName(tx, s.Name)

		if err != nil {
			slog.Error("Cannot scan for trading system '"+ s.Name +"': "+ err.Error())
			return err
		}

		if ts == nil {
			slog.Warn("Trading system '"+ s.Name +"' was not found. Skipping")
			continue
		}

		diMap, err := db.FindDailyInfoByTsIdAsMap(tx, ts.Id)

		if err != nil {
			slog.Error("Cannot retrieve DailyInfo list for trading system '"+ s.Name +"': "+ err.Error())
			return err
		}

		slog.Info("Updating trading system: "+ s.Name)

		deltaProfit := 0.0
		deltaDays   := 0
		deltaTrades := 0

		for _, inDi := range s.DailyInfo {
			if _, ok := diMap[inDi.Day]; !ok {
				di := &db.DailyInfo{
					TradingSystemId: ts.Id,
					Day            : inDi.Day,
					OpenProfit     : inDi.OpenProfit,
					ClosedProfit   : inDi.ClosedProfit,
					Position       : inDi.Position,
					NumTrades      : inDi.NumTrades,
				}

				//--- Handle the case when a trade is closed and reopened in the same
				//--- bar and in the same direction

				if di.ClosedProfit != 0 && di.NumTrades == 0 {
					di.NumTrades = 1
				}

				_ = db.AddDailyInfo(tx, di)

				//--- Add entry to map to avoid duplicates
				diMap[inDi.Day] = *di
				deltaProfit += inDi.ClosedProfit
				deltaDays   += inDi.NumTrades
				deltaTrades++

				//--- Update information on trading system

				if ts.FirstUpdate == 0 || ts.FirstUpdate > inDi.Day {
					ts.FirstUpdate = inDi.Day
				}

				if ts.LastUpdate < inDi.Day {
					ts.LastUpdate = inDi.Day
				}
			}
		}

		ts.ClosedProfit += deltaProfit
		ts.TradingDays  += deltaDays
		ts.NumTrades    += deltaTrades
		db.UpdateTradingSystem(tx, ts)
	}

	return nil
}

//=============================================================================
