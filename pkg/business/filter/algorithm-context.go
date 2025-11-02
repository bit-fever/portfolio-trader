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
	"github.com/tradalia/portfolio-trader/pkg/business/filter/algorithm/optimization"
	"github.com/tradalia/portfolio-trader/pkg/db"
	"log/slog"
)

//=============================================================================

type OptimizationContext struct {
	op *OptimizationProcess
}

//=============================================================================
//=== Constructor
//=============================================================================

func NewContext(op *OptimizationProcess) optimization.Context {
	return &OptimizationContext{
		op: op,
	}
}

//=============================================================================
//=== Methods
//=============================================================================

func (oc *OptimizationContext) FilterConfig() *optimization.FilterConfig {
	return oc.op.optReq.FilterConfig
}

//=============================================================================

func (oc *OptimizationContext) AlgorithmConfig() *optimization.AlgorithmConfig {
	return &oc.op.optReq.Algorithm.Config
}

//=============================================================================

func (oc *OptimizationContext) IsStopping() bool {
	return oc.op.stopping
}

//=============================================================================

func (oc *OptimizationContext) Baseline() db.TradingFilter {
	return *oc.op.optReq.Baseline
}

//=============================================================================

func (oc *OptimizationContext) RunAnalysis(filter *db.TradingFilter) float64 {
	return oc.op.runAnalysis(filter)
}

//=============================================================================

func (oc *OptimizationContext) LogInfo(message string) {
	slog.Info(message, "tsId", oc.op.ts.Id, "tsName", oc.op.ts.Name)
}

//=============================================================================
