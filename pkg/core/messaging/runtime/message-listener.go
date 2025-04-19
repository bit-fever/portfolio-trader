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
	"fmt"
	"github.com/bit-fever/core/msg"
	"github.com/bit-fever/portfolio-trader/pkg/business/filter"
	"github.com/bit-fever/portfolio-trader/pkg/core"
	"github.com/bit-fever/portfolio-trader/pkg/core/tradingsystem"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"github.com/bit-fever/portfolio-trader/pkg/platform"
	"github.com/vicanso/go-charts/v2"
	"gorm.io/gorm"
	"log/slog"
	"sort"
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

func addNewTrades(tx *gorm.DB, ts *db.TradingSystem, trades *[]db.Trade, newTrades []*TradeItem) (*[]db.Trade, error) {
	list := *trades

	tradeSet := map[string]bool{}
	for _, dbt := range *trades {
		tradeSet[dbt.String()] = true
	}

	for _, tr := range newTrades {
		dbTr := toDbTrade(ts.Id, tr)
		_, exists := tradeSet[dbTr.String()]
		if exists {
			continue
		}

		tradeSet[dbTr.String()] = true
		err  := db.AddTrade(tx, dbTr)

		if err != nil {
			return nil, err
		}

		list = append(list, *dbTr)

		//--- Update information on trading system

		if ts.FirstTrade == nil || ts.FirstTrade.After(*tr.EntryDate) {
			ts.FirstTrade = tr.EntryDate
		}

		if ts.LastTrade == nil || ts.LastTrade.Before(*tr.EntryDate) {
			ts.LastTrade = tr.EntryDate
		}
	}

	//--- Sort final list as new trades could be in the past

	sort.Slice(list, func(i,j int) bool {
		return list[i].ExitDate.Before(*list[j].ExitDate)
	})

	return &list, nil
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

func updateTradingSystem(tx *gorm.DB, ts *db.TradingSystem, trades *[]db.Trade, filter *db.TradingFilter) error {
	lastDays := 90
	updateLastStats(ts, trades, lastDays)
	updateActivationStatus(ts, trades, filter)

	if err := updateChart(ts, trades, lastDays); err != nil {
		return err
	}

	//--- If we got new trades, probably we have to set an idle/broken state to running

	if ts.Status == db.TsStatusIdle || ts.Status == db.TsStatusBroken {
		idleStart := time.Now().Add(-time.Hour * 24 * time.Duration(tradingsystem.IdleDays))

		if ts.LastTrade.After(idleStart) {
			ts.Status = db.TsStatusRunning
		}
	}

	return db.UpdateTradingSystem(tx, ts)
}

//=============================================================================

func updateLastStats(ts *db.TradingSystem, trades *[]db.Trade, lastDays int) {
	grossProfit := 0.0
	netProfit   := 0.0
	numTrades   := 0

	var grossEquity   []float64
	var grossDrawdown []float64
	var netEquity     []float64
	var netDrawdown   []float64

	startDate := time.Now().Add(-time.Hour * 24 * time.Duration(lastDays))

	for _, trade := range *trades {
		if trade.ExitDate.After(startDate) {
			grossProfit += trade.GrossProfit
			netProfit   += trade.GrossProfit - 2 * float64(ts.CostPerOperation)
			numTrades++

			grossEquity   = append(grossEquity, grossProfit)
			netEquity     = append(netEquity,   netProfit)
			grossDrawdown = append(grossDrawdown, 0)
			netDrawdown   = append(netDrawdown,   0)
		}
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

func updateChart(ts *db.TradingSystem, trades *[]db.Trade, lastDays int) error {
	xAxis, values := calcEquity(trades, lastDays)

	p, err := charts.LineRender(
		[][]float64{ values, },
		charts.XAxisDataOptionFunc(xAxis, charts.FalseFlag()),
		func(opt *charts.ChartOption) {
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

	if err != nil {
		slog.Error("updateChart: Cannot generate equity chart", "id", ts.Id, "error", err)
		return err
	}

	buf, err := p.Bytes()
	if err != nil {
		slog.Error("updateChart: Cannot convert equity chart to byte array", "id", ts.Id, "error", err)
		return err
	}

	err = platform.SetEquityChart(ts.Username, ts.Id, buf)
	if err != nil {
		slog.Error("updateChart: Cannot save equity chart into storage", "id", ts.Id, "error", err)
	}

	return err
}

//=============================================================================

func flatDate(date time.Time) time.Time {
	y, m, d := date.Date()
	return time.Date(y, m, d, 0,0,0,0, date.Location())
}

//=============================================================================

func calcEquity(trades *[]db.Trade, lastDays int) ([]string, []float64) {
	startDate := flatDate(time.Now().Add(-time.Hour * 24 * time.Duration(lastDays)))
	xLabel    := startDate.Format(time.DateOnly)

	var xAxis   []string
	var values  []float64

	for i := 0; i<=lastDays; i++ {
		xAxis  = append(xAxis, xLabel)
		values = append(values, 0)

		xLabel = ""
	}

	for _, trade := range *trades {
		if trade.ExitDate.After(startDate) {
			day := trade.ExitDate.Sub(startDate) / (time.Hour * 24)
			values[day] += trade.GrossProfit
		}
	}

	for i := 1; i<len(values); i++ {
		values[i] += values[i -1]
	}

	return xAxis, values
}

//=============================================================================
