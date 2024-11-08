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

package tradingsystem

import (
	"github.com/bit-fever/portfolio-trader/pkg/app"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

//=============================================================================

var IdleDays   = 7
var BrokenDays = 14

//=============================================================================

func InitUpdaterProcess(cfg *app.Config) *time.Ticker {
	ticker := time.NewTicker(8 * time.Hour)

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
	slog.Info("TradingSystemUpdater: Starting")
	start := time.Now()

	list, err := GetTradingSystemsInIdle()
	if err != nil {
		slog.Error("Cannot get list of trading systems. Update aborted", "error", err)
	} else {
		slog.Info("TradingSystemUpdater: Processing trading systems", "count", len(*list))

		for _, ts := range *list {
			err = updateTradingSystem(&ts)
			if err != nil {
				slog.Error("Cannot update trading system", "id", ts.Id, "error", err)
			}
		}
	}

	duration := time.Now().Sub(start).Seconds()
	slog.Info("TradingSystemUpdater: Ended", "seconds", duration)
}

//=============================================================================

func GetTradingSystemsInIdle() (*[]db.TradingSystem, error){
	var list *[]db.TradingSystem
	var err error

	err = db.RunInTransaction(func (tx *gorm.DB) error {
		list, err = db.GetTradingSystemsInIdle(tx, IdleDays)
		return err
	})

	return list,err
}

//=============================================================================

func updateTradingSystem(ts *db.TradingSystem) error {
	if ts.Status == db.TsStatusRunning {
		ts.Status = db.TsStatusIdle
	} else if ts.Status == db.TsStatusIdle {
		brokenDate := time.Now().Add(-time.Hour * 24 * time.Duration(BrokenDays))

		if ts.LastTrade.Before(brokenDate) {
			ts.Status = db.TsStatusBroken
		}
	} else {
		return nil
	}

	return db.RunInTransaction(func (tx *gorm.DB) error {
		return db.UpdateTradingSystem(tx, ts)
	})
}

//=============================================================================
