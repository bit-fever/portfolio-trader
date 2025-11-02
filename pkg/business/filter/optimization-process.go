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

package filter

import (
	"github.com/tradalia/portfolio-trader/pkg/business/filter/algorithm"
	"github.com/tradalia/portfolio-trader/pkg/business/filter/algorithm/optimization"
	"github.com/tradalia/portfolio-trader/pkg/db"
	"log/slog"
	"math/rand"
	"time"
)

//=============================================================================

const MaxResultSize = 1000

//=============================================================================
//===
//=== OptimizationProcess
//===
//=============================================================================

type OptimizationProcess struct {
	ts              *db.TradingSystem
	trades          *[]db.Trade
	optReq          *OptimizationRequest
	info            *OptimizationInfo
	fitnessFunction FitnessFunction
	stopping        bool
}

//=============================================================================

func (op *OptimizationProcess) Start() {
	field     := op.optReq.FieldToOptimize
	startDate := op.optReq.StartDate

	op.fitnessFunction = GetFitnessFunction(field)

	algo := algorithm.New(op.optReq.Algorithm.Type)
	ctx  := NewContext(op)
	algo.Init(ctx)

	fc := op.optReq.FilterConfig
	op.info = NewOptimizationInfo(MaxResultSize, field, fc, algo.StepsCount(), op.calcBaseValue(), startDate)

	go op.generate(algo)
}

//=============================================================================

func (op *OptimizationProcess) Stop() {
	slog.Info("Stop: Stopping optimization process", "tsId", op.ts.Id)
	op.stopping = true
}

//=============================================================================

func (op *OptimizationProcess) GetInfo() *OptimizationInfo {
	return op.info
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

//--- GoRoutine

func (op *OptimizationProcess) generate(algo optimization.Algorithm) {
	slog.Info("generate: Started", "tsId", op.ts.Id, "tsName", op.ts.Name, "algorithm", op.optReq.Algorithm)

	algo.Optimize()

	for ; !op.info.isStatusComplete(); {
		time.Sleep(time.Second * 1)
	}

	slog.Info("generate: Complete.")
}

//=============================================================================

func (op *OptimizationProcess) runAnalysis(filter *db.TradingFilter) float64 {
	sum := RunAnalysis(op.ts, filter, op.trades).Summary
	run := op.createRun(filter, &sum)
	op.info.addResult(run)

	return run.FitnessValue
}

//=============================================================================

func (op *OptimizationProcess) createRun(filter *db.TradingFilter, sum *Summary) *Run {
	r := &Run{
		Filter      : filter,
		NetProfit   : sum.FilProfit,
		AvgTrade    : sum.FilAverageTrade,
		MaxDrawdown : sum.FilMaxDrawdown,
		random      : rand.Int(),
	}

	r.FitnessValue = op.fitnessFunction(r)

	return r
}

//=============================================================================

func (op *OptimizationProcess) calcBaseValue() float64 {
	baseline := op.optReq.Baseline

	sum := RunAnalysis(op.ts, baseline, op.trades).Summary
	run := op.createRun(baseline, &sum)

	return run.FitnessValue
}

//=============================================================================
