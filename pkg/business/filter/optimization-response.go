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
	"time"
)

//=============================================================================
//===
//=== OptimizationResponse
//===
//=============================================================================

type OptimizationResponse struct {
	CurrStep        uint          `json:"currStep"`
	MaxSteps        uint          `json:"maxSteps"`
	StartTime       time.Time     `json:"startTime"`
	EndTime         time.Time     `json:"endTime"`
	Status          string        `json:"status"`
	Runs            []any         `json:"runs"`
	BaseValue       float64       `json:"baseValue"`
	BestValue       float64       `json:"bestValue"`
	FieldToOptimize string        `json:"fieldToOptimize"`
	Duration        int64         `json:"duration"`
	Filters struct {
		PosProfit bool `json:"posProfit"`
		OldVsNew  bool `json:"oldVsNew"`
		WinPerc   bool `json:"winPerc"`
		EquVsAvg  bool `json:"equVsAvg"`
	} `json:"filters"`
}

//=============================================================================

func NewOptimizationResponse(info *OptimizationInfo) *OptimizationResponse {
	or := &OptimizationResponse{}
	or.CurrStep  = info.CurrStep
	or.MaxSteps  = info.MaxSteps
	or.StartTime = info.StartTime
	or.EndTime   = info.EndTime
	or.Status    = info.Status
	or.BaseValue = info.BaseValue
	or.BestValue = info.BestValue

	or.FieldToOptimize   = info.FieldToOptimize
	or.Filters.PosProfit = info.Filters.PosProfit
	or.Filters.OldVsNew  = info.Filters.OldVsNew
	or.Filters.WinPerc   = info.Filters.WinPerc
	or.Filters.EquVsAvg  = info.Filters.EquVsAvg

	or.Runs     = info.GetRuns()
	or.Duration = int64(time.Now().Sub(info.StartTime).Seconds())

	return or
}

//=============================================================================
