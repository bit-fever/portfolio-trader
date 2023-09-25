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
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"github.com/bit-fever/portfolio-trader/pkg/model/config"
	"github.com/bit-fever/portfolio-trader/pkg/tool"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

//=============================================================================

func StartPeriodicScan(cfg *config.Config) *time.Ticker {

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

func run(cfg *config.Config) {
	url := cfg.Scan.Address
	log.Println("Starting to fetch strategies from: "+ url)

	var data []strategy

	err := tool.DoGet(tool.GetClient("ws"), url, &data)

	if err == nil {
		log.Println("Got "+ strconv.Itoa(len(data))+ " strategies")

		db.RunInTransaction(func (tx *gorm.DB) error {
			return updateDb(tx, data)
		})
	} else {
		log.Println("Cannot connect to url. Error is: "+ err.Error())
	}

	log.Println("Ending to fetch strategies")
}

//=============================================================================

func updateDb(tx *gorm.DB, inStrategies []strategy) error {
	log.Println("Updating database...")

	p, err := db.GetOrCreatePortfolio(tx, "Main", &db.Portfolio{
		Name: "Main",
	})

	if err != nil {
		log.Println("Cannot add the first portfolio: "+ err.Error())
		return err
	}

	for _, s := range inStrategies {
		log.Println("Updating trading system: "+ s.Name)

		in, err := db.GetOrCreateInstrument(tx, s.Ticker, &db.Instrument{
			Ticker: s.Ticker,
			Name: s.Ticker,
		})

		if err != nil {
			log.Println("Cannot add the instrument with ticker '"+ s.Ticker +"': "+ err.Error())
			return err
		}

		ts, err := db.GetOrCreateTradingSystem(tx, s.Name, &db.TradingSystem{
			Code: s.Name,
			Name: s.Name,
			InstrumentId: in.Id,
			PortfolioId: p.Id,
			FilterType: 0,
			SuggestedAction: 0,
		})

		if err != nil {
			log.Println("Cannot add the trading system '"+ s.Name +"': "+ err.Error())
			return err
		}

		diMap, err := db.FindDailyInfoByTsIdAsMap(tx, ts.Id)

		if err != nil {
			log.Println("Cannot retrieve DailyInfo list for trading system '"+ s.Name +"': "+ err.Error())
			return err
		}

		deltaTrades := 0

		for _, inDi := range s.DailyInfo {
			if _, ok := diMap[inDi.Day]; !ok {
				_ = db.AddTsDailyInfo(tx, &db.TsDailyInfo{
					TradingSystemId: ts.Id,
					Day: inDi.Day,
					OpenProfit: inDi.OpenProfit,
					Position: inDi.Position,
					NumTrades: inDi.NumTrades,
				})

				deltaTrades++

				//--- Update information on trading system

				if ts.FirstUpdate == 0 || ts.FirstUpdate > inDi.Day {
					ts.FirstUpdate = inDi.Day
				}

				if ts.LastUpdate < inDi.Day {
					ts.LastUpdate = inDi.Day
					ts.LastPl     = inDi.OpenProfit
					ts.NumTrades  = inDi.NumTrades
				}
			}
		}

		ts.TradingDays += deltaTrades
		tx.Updates(ts)
	}

	return nil
}

//=============================================================================
