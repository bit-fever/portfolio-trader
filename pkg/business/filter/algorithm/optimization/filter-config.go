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

package optimization

import (
	"errors"
	"math/rand"
	"strconv"
)

//=============================================================================
//===
//=== FilterConfig
//===
//=============================================================================

const MaxTradesLength     = 300
const MaxTrendlineValue   = 200
const MaxOldNewPercentage = 200
const MaxWinningPercentage= 100

//=============================================================================

type FilterConfig struct {
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

func (fc *FilterConfig) Validate() error {
	if err := fc.PosProLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := fc.OldNewOldLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := fc.OldNewNewLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := fc.OldNewOldPerc.Validate(1, MaxOldNewPercentage); err != nil {
		return err
	}

	if err := fc.WinPercLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := fc.WinPercPerc.Validate(1, MaxWinningPercentage); err != nil {
		return err
	}

	if err := fc.EquAvgLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := fc.TrendlineLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := fc.TrendlineValue.Validate(1, MaxTrendlineValue); err != nil {
		return err
	}

	return nil
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

	//--- Caching
	steps    *[]int
}

//=============================================================================

func (f *FieldOptimization) StepsCount() uint {
	if !f.Enabled {
		return 1
	}

	return uint((f.MaxValue - f.MinValue) / f.Step) +1
}

//=============================================================================

func (f *FieldOptimization) Steps() *[]int {
	if f.steps != nil {
		return f.steps
	}

	list := []int{}

	if !f.Enabled {
		list = append(list, f.CurValue)
	} else {
		for i := f.MinValue; i<= f.MaxValue; i += f.Step {
			list = append(list, i)
		}
	}

	f.steps = &list

	return &list
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

func (f *FieldOptimization) RandomValue() int {
	list := f.Steps()
	idx  := rand.Intn(len(*list))

	return (*list)[idx]
}

//=============================================================================
