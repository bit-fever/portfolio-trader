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

import "errors"

//=============================================================================
//===
//=== OptimizationRequest
//===
//=============================================================================

const FieldToOptimizeNetProfit            = "net-profit"
const FieldToOptimizeAvgTrade             = "avg-trade"
const FieldToOptimizeDrawDown             = "drawdown"
const FieldToOptimizeNetProfitAndAvgTrade = "net-profit*avg-trade"

//=============================================================================

type OptimizationRequest struct {
	FieldToOptimize string `json:"fieldToOptimize"       binding:"required"`
	EnablePosProfit bool   `json:"enablePosProfit"`
	EnableOldNew    bool   `json:"enableOldNew"`
	EnableWinPerc   bool   `json:"enableWinPerc"`
	EnableEquAvg    bool   `json:"enableEquAvg"`
	CombineFilters  bool   `json:"combineFilters"`

	PosProLen     FieldOptimization  `json:"posProLen"`
	OldNewOldLen  FieldOptimization  `json:"oldNewOldLen"`
	OldNewNewLen  FieldOptimization  `json:"oldNewNewLen"`
	OldNewOldPerc FieldOptimization  `json:"oldNewOldPerc"`
	WinPercLen    FieldOptimization  `json:"winPercLen"`
	WinPercPerc   FieldOptimization  `json:"winPercPerc"`
	EquAvgLen     FieldOptimization  `json:"equAvgLen"`
}

//=============================================================================

func (f *OptimizationRequest) StepsCount() uint {
	if f.CombineFilters {
		return	f.stepsCountPosProfit(1) *
				f.stepsCountOldNew   (1) *
				f.stepsCountWinPerc  (1) *
				f.stepsCountEquAvg   (1)
	} else {
		return	f.stepsCountPosProfit(0) +
				f.stepsCountOldNew   (0) +
				f.stepsCountWinPerc  (0) +
				f.stepsCountEquAvg   (0)
	}
}
//=============================================================================

func (f *OptimizationRequest) Validate() error {
	if  f.FieldToOptimize != FieldToOptimizeNetProfit &&
		f.FieldToOptimize != FieldToOptimizeAvgTrade &&
		f.FieldToOptimize != FieldToOptimizeDrawDown &&
		f.FieldToOptimize != FieldToOptimizeNetProfitAndAvgTrade {
		return errors.New("Invalid field to optimize: "+ f.FieldToOptimize)
	}

	//TODO: validate the other params here

	return nil
}

//=============================================================================

func (f *OptimizationRequest) stepsCountPosProfit(zero uint) uint {
	if f.EnablePosProfit {
		return f.PosProLen.StepsCount()
	}

	return zero
}

//=============================================================================

func (f *OptimizationRequest) stepsCountOldNew(zero uint) uint {
	if f.EnableOldNew {
		return f.OldNewOldLen.StepsCount() * f.OldNewNewLen.StepsCount() * f.OldNewOldPerc.StepsCount()
	}

	return zero
}

//=============================================================================

func (f *OptimizationRequest) stepsCountWinPerc(zero uint) uint {
	if f.EnableWinPerc {
		return f.WinPercLen.StepsCount() * f.WinPercPerc.StepsCount()
	}

	return zero
}

//=============================================================================

func (f *OptimizationRequest) stepsCountEquAvg(zero uint) uint {
	if f.EnableEquAvg {
		return f.EquAvgLen.StepsCount()
	}

	return zero
}

//=============================================================================
//===
//=== FieldOptimization
//===
//=============================================================================

type FieldOptimization struct {
	Enabled  bool  `json:"enabled"`
	CurValue int   `json:"curValue"   binding:"min=1,max=600"`
	MinValue int   `json:"minValue"   binding:"min=1,max=600"`
	MaxValue int   `json:"maxValue"   binding:"min=1,max=600"`
	Step     int   `json:"step"       binding:"min=1,max=100"`
}

//=============================================================================

func (f *FieldOptimization) StepsCount() uint {
	if !f.Enabled {
		return 1
	}

	return uint((f.MaxValue - f.MinValue) / f.Step) +1
}

//=============================================================================

func (f *FieldOptimization) Steps() []int {
	list := []int{}

	if !f.Enabled {
		list = append(list, f.CurValue)
	} else {
		for i := f.MinValue; i<= f.MaxValue; i += f.Step {
			list = append(list, i)
		}
	}

	return list
}

//=============================================================================
