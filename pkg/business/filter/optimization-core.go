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
	"github.com/bit-fever/portfolio-trader/pkg/core"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"log/slog"
	"sync"
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
	num := 8
	workers.Init(num, 100)
}

//=============================================================================
//===
//=== API methods
//===
//=============================================================================

func StartOptimization(ts *db.TradingSystem, data *[]db.DailyInfo, params *OptimizationRequest) error {
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
		data  : data,
		params: params,
	}

	err := fop.Start()

	if err == nil {
		jobs.m[ts.Id] = fop
	}

	return err
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
