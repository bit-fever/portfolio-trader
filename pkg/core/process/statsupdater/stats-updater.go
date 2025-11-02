//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

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

package statsupdater

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/tradalia/portfolio-trader/pkg/app"
	"github.com/tradalia/portfolio-trader/pkg/core"
	"github.com/tradalia/portfolio-trader/pkg/db"
	"github.com/tradalia/portfolio-trader/pkg/platform"
	"github.com/vicanso/go-charts/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
	"gorm.io/gorm"
)

//=============================================================================

const (
	ChartTypeTime   = "time"
	ChartTypeTrades = "trades"
)

const LastDays = 180

//=============================================================================

func Init(cfg *app.Config) *time.Ticker {
	ticker := time.NewTicker(12 * time.Hour)

	go func() {
		//--- Wait 6 secs to allow the system to boot properly
		time.Sleep(6 * time.Second)
		run(cfg)

		for range ticker.C {
			run(cfg)
		}
	}()

	return ticker
}

//=============================================================================

func run(cfg *app.Config) {
	slog.Info("StatsUpdater: Starting")
	start := time.Now()

	users, err := getUsersWithTradingSystems()
	if err != nil {
		slog.Error("StatsUpdater: Cannot get list of users with trading systems. Update aborted", "error", err)
	} else {
		for _, user := range users {
			updateSystemsForUser(user)
		}
	}

	duration := time.Now().Sub(start).Seconds()
	slog.Info("StatsUpdater: Ended", "seconds", duration)
}

//=============================================================================

func getUsersWithTradingSystems() ([]string, error){
	var list []string
	var err error

	err = db.RunInTransaction(func (tx *gorm.DB) error {
		list, err = db.GetUsersWithTradingSystems(tx)
		return err
	})

	return list,err
}

//=============================================================================

func updateSystemsForUser(user string) {
	list, err := getTradingSystemsByUser(user)
	if err != nil {
		slog.Error("updateSystemsForUser: Cannot get list of trading systems. Update aborted for this user", "user", user, "error", err)
	} else {
		slog.Info("updateSystemsForUser: Updating trading systems", "user", user, "count", len(*list))

		for _, ts := range *list {
			updateTradingSystem(&ts)
		}
	}
}

//=============================================================================

func getTradingSystemsByUser(user string) (*[]db.TradingSystem, error){
	var list *[]db.TradingSystem
	var err error

	err = db.RunInTransaction(func (tx *gorm.DB) error {
		list, err = db.GetTradingSystemsByUser(tx,  user)
		return err
	})

	return list,err
}

//=============================================================================

func updateTradingSystem(ts *db.TradingSystem) {
	lastDays := LastDays
	trades, err := getTradingSystemTrades(ts.Id, lastDays)
	if err != nil {
		slog.Error("updateTradingSystem: Cannot get the list of trades for trading system. Skipping", "user", ts.Username, "id", ts.Id, "error", err)
	} else {
		updateLastStats(ts, trades)
		err = updateChart(ts, trades, lastDays)
		if err == nil {
			err = db.RunInTransaction(func (tx *gorm.DB) error {
				return db.UpdateTradingSystem(tx, ts)
			})
			if err != nil {
				slog.Error("updateTradingSystem: Cannot update trading system. Skipping", "user", ts.Username, "id", ts.Id, "error", err)
			}
		}
	}
}

//=============================================================================

func getTradingSystemTrades(id uint, lastDays int) (*[]db.Trade, error) {
	var list *[]db.Trade
	var err error

	fromDate := time.Now().Add(-time.Hour * 24 * time.Duration(lastDays))

	err = db.RunInTransaction(func (tx *gorm.DB) error {
		list, err = db.FindTradesByTsIdFromTime(tx, id, &fromDate, nil)
		return err
	})

	return list,err
}

//=============================================================================

func updateLastStats(ts *db.TradingSystem, trades *[]db.Trade) {
	grossProfit := 0.0
	netProfit   := 0.0
	numTrades   := 0

	var grossEquity   []float64
	//var grossDrawdown []float64
	var netEquity     []float64
	//var netDrawdown   []float64

	for _, trade := range *trades {
		grossProfit += trade.GrossProfit
		netProfit   += trade.GrossProfit - 2 * ts.CostPerOperation
		numTrades++

		grossEquity   = append(grossEquity, grossProfit)
		netEquity     = append(netEquity,   netProfit)
		//grossDrawdown = append(grossDrawdown, 0)
		//netDrawdown   = append(netDrawdown,   0)
	}

	//maxGrossDD := core.CalcDrawDown(&grossEquity, &grossDrawdown)
	//maxNetDD   := core.CalcDrawDown(&netEquity,   &netDrawdown)

	ts.LastNetProfit   = core.Trunc2d(netProfit)
	ts.LastNumTrades   = numTrades
	ts.LastNetAvgTrade = 0

	if numTrades != 0 {
		ts.LastNetAvgTrade = core.Trunc2d(netProfit / float64(numTrades))
	}
}

//=============================================================================

func updateChart(ts *db.TradingSystem, trades *[]db.Trade, lastDays int) error {
	er := platform.NewEquityRequest()
	er.Username = ts.Username

	//--- Time based image

	p, err := buildEquityChartTime(ts, trades, lastDays)
	if err != nil {
		slog.Error("updateChart: Cannot generate equity chart (time)", "id", ts.Id, "error", err)
		return err
	}

	buf, err := p.Bytes()
	if err != nil {
		slog.Error("updateChart: Cannot convert equity chart (time) to byte array", "id", ts.Id, "error", err)
		return err
	}

	er.Images[ChartTypeTime] = buf

	//--- Trades based image

	p, err = buildEquityChartTrades(ts, trades)
	if err != nil {
		slog.Error("updateChart: Cannot generate equity chart (trades)", "id", ts.Id, "error", err)
		return err
	}

	//--- If the TS is new, we don't have any trade and need to skip the chart generation

	if p == nil {
		return nil
	}

	buf, err = p.Bytes()
	if err != nil {
		slog.Error("updateChart: Cannot convert equity chart (trades) to byte array", "id", ts.Id, "error", err)
		return err
	}

	er.Images[ChartTypeTrades] = buf

	//--- Send images to storage manager

	err = platform.SetEquityChart(ts.Id, er)
	if err != nil {
		slog.Error("updateChart: Cannot save equity chart into storage", "id", ts.Id, "error", err)
	}

	return err
}

//=============================================================================

func buildEquityChartTime(ts *db.TradingSystem, trades *[]db.Trade, lastDays int) (*charts.Painter, error) {
	xAxis, values := calcEquityTime(trades, lastDays, float64(ts.CostPerOperation))

	return charts.LineRender(
		[][]float64{ values, },
		charts.XAxisDataOptionFunc(xAxis, charts.FalseFlag()),
		func(opt *charts.ChartOption) {
			opt.BackgroundColor = drawing.ColorFromHex("F6F9FF")
			opt.XAxis.SplitNumber = 30
			opt.XAxis.FontSize = 8
			opt.SymbolShow = charts.FalseFlag()
			opt.LineStrokeWidth = 2
			opt.ValueFormatter = func(f float64) string {
				return fmt.Sprintf("%.0f", f)
			}
			opt.Width  = 400
			opt.Height = 200
			opt.YAxisOptions = []charts.YAxisOption{ { FontSize: 8 }}
			opt.Padding      = charts.Box{ Top: 8, Left: 4, Right: 4, Bottom: 0}
		},
	)
}

//=============================================================================

func calcEquityTime(trades *[]db.Trade, lastDays int, costPerOper float64) ([]string, []float64) {
	startDate := time.Now().Add(-time.Hour * 24 * time.Duration(lastDays))
	startDate  = flatDate(&startDate)
	xLabel    := startDate.Format(time.DateOnly)

	var xAxis  []string
	var values []float64

	for i := 0; i<=lastDays; i++ {
		xAxis  = append(xAxis, xLabel)
		values = append(values, 0)

		xLabel = ""
	}

	for _, trade := range *trades {
		day := trade.ExitDate.Sub(startDate) / (time.Hour * 24)
		values[day] += trade.GrossProfit - 2 * costPerOper
	}

	for i := 1; i<len(values); i++ {
		values[i] += values[i -1]
	}

	return xAxis, values
}

//=============================================================================

func buildEquityChartTrades(ts *db.TradingSystem, trades *[]db.Trade) (*charts.Painter, error) {
	xAxis, values := calcEquityTrades(trades, float64(ts.CostPerOperation))

	if xAxis == nil {
		return nil, nil
	}

	return charts.LineRender(
		[][]float64{ values, },
		charts.XAxisDataOptionFunc(xAxis, charts.FalseFlag()),
		func(opt *charts.ChartOption) {
			opt.BackgroundColor = drawing.ColorFromHex("F6F9FF")
			opt.XAxis.SplitNumber = 5
			opt.XAxis.FontSize = 8
			opt.SymbolShow = charts.FalseFlag()
			opt.LineStrokeWidth = 2
			opt.ValueFormatter = func(f float64) string {
				return fmt.Sprintf("%.0f", f)
			}
			opt.Width  = 400
			opt.Height = 200
			opt.YAxisOptions = []charts.YAxisOption{ { FontSize: 8 }}
			opt.Padding      = charts.Box{ Top: 8, Left: 4, Right: 4, Bottom: 0}
		},
	)
}

//=============================================================================

func calcEquityTrades(trades *[]db.Trade, costPerOper float64) ([]string, []float64) {
	var xAxis  []string
	var values []float64

	netProfit := 0.0

	for i, trade := range *trades {
		netProfit += trade.GrossProfit - 2 * costPerOper

		xAxis  = append(xAxis,  strconv.Itoa(i +1))
		values = append(values, netProfit)
	}

	return xAxis, values
}

//=============================================================================

func flatDate(date *time.Time) time.Time {
	y, m, d := date.Date()
	return time.Date(y, m, d, 0,0,0,0, date.Location())
}

//=============================================================================
