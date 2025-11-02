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

package filter

import (
	"github.com/tradalia/portfolio-trader/pkg/core"
	"github.com/tradalia/portfolio-trader/pkg/db"
	"log/slog"
	"sync"
	"time"
)

//=============================================================================

var jobs = struct {
	sync.RWMutex
	m map[uint]*OptimizationProcess
}{m: make(map[uint]*OptimizationProcess)}

//-----------------------------------------------------------------------------

var workers = core.WorkerPool{}

//=============================================================================
//===
//=== Init
//===
//=============================================================================

func init() {
	num := 1
	workers.Init(num, 100)
	go periodicCleanup()
}

//=============================================================================
//===
//=== API methods
//===
//=============================================================================

func StartOptimization(ts *db.TradingSystem, trades *[]db.Trade, or *OptimizationRequest) {
	jobs.Lock()
	defer jobs.Unlock()

	fop, ok := jobs.m[ts.Id]
	if ok {
		slog.Error("Stopping a previous optimization process", "tsId", ts.Id)
		fop.Stop()
		delete(jobs.m, ts.Id)
	}

	fop = &OptimizationProcess{
		ts    : ts,
		trades: trades,
		optReq: or,
	}

	fop.Start()
	jobs.m[ts.Id] = fop
}

//=============================================================================

func StopOptimization(tsId uint) error {
	jobs.Lock()
	defer jobs.Unlock()

	fop, ok := jobs.m[tsId]
	if ok {
		fop.Stop()
	}

	return nil
}

//=============================================================================

func GetOptimizationInfo(tsId uint) *OptimizationInfo {
	jobs.Lock()
	defer jobs.Unlock()

	fop, ok := jobs.m[tsId]
	if !ok {
		return &OptimizationInfo{
			Status: OptimStatusIdle,
		}
	}

	return fop.GetInfo()
}

//=============================================================================
//===
//=== Cleanup process
//===
//=============================================================================

func periodicCleanup() {
	for {
		time.Sleep(time.Minute * 5)
		purge()
	}
}

//=============================================================================

func purge() {
	jobs.Lock()
	defer jobs.Unlock()

	for tsId, op := range jobs.m {
		if op.info.Status == OptimStatusComplete {
			delta := time.Now().Sub(op.info.EndTime)
			if delta.Minutes() >= 30 {
				slog.Info("purge: Purging optimization process entry for trading system", "tsId", tsId)
				delete(jobs.m, tsId)
			}
		}
	}
}

//=============================================================================
