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
	"strconv"
)

//=============================================================================
//===
//=== OptimizationRequest
//===
//=============================================================================

const AlgorithmSimple  = "simple"
const AlgorithmGenetic = "genetic"

const FieldToOptimizeNetProfit            = "net-profit"
const FieldToOptimizeAvgTrade             = "avg-trade"
const FieldToOptimizeDrawDown             = "drawdown"
const FieldToOptimizeNetProfitAndAvgTrade = "net-profit*avg-trade"

const MaxTradesLength     = 300
const MaxTrendlineValue   = 200
const MaxOldNewPercentage = 200
const MaxWinningPercentage= 100

//=============================================================================

type OptimizationRequest struct {
	FieldToOptimize string `json:"fieldToOptimize"`
	Algorithm       string `json:"algorithm"`
	EnablePosProfit bool   `json:"enablePosProfit"`
	EnableOldNew    bool   `json:"enableOldNew"`
	EnableWinPerc   bool   `json:"enableWinPerc"`
	EnableEquAvg    bool   `json:"enableEquAvg"`
	EnableTrendline bool   `json:"enableTrendline"`

	PosProLen      FieldOptimization  `json:"posProLen"`
	OldNewOldLen   FieldOptimization  `json:"oldNewOldLen"`
	OldNewNewLen   FieldOptimization  `json:"oldNewNewLen"`
	OldNewOldPerc  FieldOptimization  `json:"oldNewOldPerc"`
	WinPercLen     FieldOptimization  `json:"winPercLen"`
	WinPercPerc    FieldOptimization  `json:"winPercPerc"`
	EquAvgLen      FieldOptimization  `json:"equAvgLen"`
	TrendlineLen   FieldOptimization  `json:"trendlineLen"`
	TrendlineValue FieldOptimization  `json:"trendlineValue"`
}

//=============================================================================

func (r *OptimizationRequest) StepsCount() uint {
	switch r.Algorithm {
		case AlgorithmSimple:
			return	r.stepsCountPosProfit(0) +
					r.stepsCountOldNew   (0) +
					r.stepsCountWinPerc  (0) +
					r.stepsCountEquAvg   (0) +
					r.stepsCountTrendline(0)

		case AlgorithmGenetic:
			return 0

		//case FullScan:
		//	return	r.stepsCountPosProfit(1) *
		//			r.stepsCountOldNew   (1) *
		//			r.stepsCountWinPerc  (1) *
		//			r.stepsCountEquAvg   (1) *
		//			r.stepsCountTrendline(1)

		default:
			panic("Unknown algorithm: "+ r.Algorithm)
	}
}

//=============================================================================

func (r *OptimizationRequest) Validate() error {
	if  r.FieldToOptimize != FieldToOptimizeNetProfit &&
		r.FieldToOptimize != FieldToOptimizeAvgTrade &&
		r.FieldToOptimize != FieldToOptimizeDrawDown &&
		r.FieldToOptimize != FieldToOptimizeNetProfitAndAvgTrade {
		return errors.New("Invalid field to optimize: "+ r.FieldToOptimize)
	}

	if  r.Algorithm != AlgorithmSimple &&
		r.Algorithm != AlgorithmGenetic {
		return errors.New("Invalid optimization algorithm: "+ r.Algorithm)
	}

	//--- Validate fields

	if err := r.PosProLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := r.OldNewOldLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := r.OldNewNewLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := r.OldNewOldPerc.Validate(1, MaxOldNewPercentage); err != nil {
		return err
	}

	if err := r.WinPercLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := r.WinPercPerc.Validate(1, MaxWinningPercentage); err != nil {
		return err
	}

	if err := r.EquAvgLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := r.TrendlineLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := r.TrendlineValue.Validate(1, MaxTrendlineValue); err != nil {
		return err
	}

	return nil
}

//=============================================================================

func (r *OptimizationRequest) stepsCountPosProfit(zero uint) uint {
	if r.EnablePosProfit {
		return r.PosProLen.StepsCount()
	}

	return zero
}

//=============================================================================

func (r *OptimizationRequest) stepsCountOldNew(zero uint) uint {
	if r.EnableOldNew {
		return r.OldNewOldLen.StepsCount() * r.OldNewNewLen.StepsCount() * r.OldNewOldPerc.StepsCount()
	}

	return zero
}

//=============================================================================

func (r *OptimizationRequest) stepsCountWinPerc(zero uint) uint {
	if r.EnableWinPerc {
		return r.WinPercLen.StepsCount() * r.WinPercPerc.StepsCount()
	}

	return zero
}

//=============================================================================

func (r *OptimizationRequest) stepsCountEquAvg(zero uint) uint {
	if r.EnableEquAvg {
		return r.EquAvgLen.StepsCount()
	}

	return zero
}

//=============================================================================

func (r *OptimizationRequest) stepsCountTrendline(zero uint) uint {
	if r.EnableTrendline {
		return r.TrendlineLen.StepsCount() * r.TrendlineValue.StepsCount()
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
	CurValue int   `json:"curValue"`
	MinValue int   `json:"minValue"`
	MaxValue int   `json:"maxValue"`
	Step     int   `json:"step"`
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

func (f *FieldOptimization) Validate(min, max int) error {
	if f.Enabled {
		if f.MinValue < min || f.MinValue > max {
			return errors.New("min value out of range ["+ strconv.Itoa(min) +".."+ strconv.Itoa(max) +"]")
		}

		if f.MaxValue < min || f.MaxValue > max {
			return errors.New("max value out of range ["+ strconv.Itoa(min) +".."+ strconv.Itoa(max) +"]")
		}

		if f.MinValue > f.MaxValue {
			return errors.New("min value greater than max value")
		}

		if f.Step < 1 || f.Step > max {
			return errors.New("step value out of range [1.."+ strconv.Itoa(max) +"]")
		}
	} else {
		if f.CurValue < min || f.CurValue > max {
			return errors.New("current value out of range ["+ strconv.Itoa(min) +".."+ strconv.Itoa(max) +"]")
		}
	}

	return nil
}

//=============================================================================
