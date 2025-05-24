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
	"errors"
	"github.com/bit-fever/portfolio-trader/pkg/business/filter/algorithm"
	"github.com/bit-fever/portfolio-trader/pkg/business/filter/algorithm/optimization"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"time"
)

//=============================================================================
//===
//=== OptimizationRequest
//===
//=============================================================================

const FieldToOptimizeNetProfit              = "netProfit"
const FieldToOptimizeAvgTrade               = "avgTrade"
const FieldToOptimizeDrawDown               = "maxDD"
const FieldToOptimizeNetProfitAvgTrade      = "netProfit*avgTrade"
const FieldToOptimizeNetProfitAvgTradeMaxDD = "netProfit*avgTrade/maxDD"

//=============================================================================

type AlgorithmSpec struct {
	Type   string                       `json:"type"`
	Config optimization.AlgorithmConfig `json:"config"`
}

//=============================================================================

type OptimizationRequest struct {
	StartDate       *time.Time                 `json:"startDate,omitempty"`
	FieldToOptimize string                     `json:"fieldToOptimize"`
	FilterConfig    *optimization.FilterConfig `json:"filterConfig"`
	Algorithm       *AlgorithmSpec             `json:"algorithm"`
	Baseline        *db.TradingFilter          `json:"baseline"`
}

//=============================================================================

func (r *OptimizationRequest) Validate() error {
	if  r.FieldToOptimize != FieldToOptimizeNetProfit         &&
		r.FieldToOptimize != FieldToOptimizeAvgTrade          &&
		r.FieldToOptimize != FieldToOptimizeDrawDown          &&
		r.FieldToOptimize != FieldToOptimizeNetProfitAvgTrade &&
		r.FieldToOptimize != FieldToOptimizeNetProfitAvgTradeMaxDD {
		return errors.New("Invalid field to optimize: "+ r.FieldToOptimize)
	}

	algoType := r.Algorithm.Type

	if  algoType != algorithm.Simple && algoType != algorithm.Genetic {
		return errors.New("Invalid optimization algorithm: "+ algoType)
	}

	if err := r.FilterConfig.Validate(); err != nil {
		return err
	}

	return nil
}

//=============================================================================
