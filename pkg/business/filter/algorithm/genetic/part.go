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

package genetic

import (
	"github.com/tradalia/portfolio-trader/pkg/business/filter/algorithm/optimization"
	"github.com/tradalia/portfolio-trader/pkg/db"
)

//=============================================================================

type Part interface {
	Mutate()
	Apply(filter *db.TradingFilter)
}

//=============================================================================
//===
//=== PosProfitPart
//===
//=============================================================================

type PosProfitPart struct {
	lenEnabled bool
	len        int

	lenFo *optimization.FieldOptimization
}

//=============================================================================

func NewPosProfitPart(fc *optimization.FilterConfig) Part {
	p := &PosProfitPart{
		lenEnabled: fc.PosProLen.Enabled,
		len       : fc.PosProLen.CurValue,
		lenFo     : &fc.PosProLen,
	}

	if p.lenEnabled {
		p.len = p.lenFo.RandomValue()
	}

	return p
}

//=============================================================================

func (p *PosProfitPart) Mutate() {

}

//=============================================================================

func (p *PosProfitPart) Apply(f *db.TradingFilter) {

}

//=============================================================================
//===
//=== EquityVsAvgPart
//===
//=============================================================================

type EquityVsAvgPart struct {
	lenEnabled bool
	len        int

	lenFo *optimization.FieldOptimization
}

//=============================================================================

func NewEquityVsAvgPart(fc *optimization.FilterConfig) Part {
	p := &EquityVsAvgPart{
		lenEnabled: fc.EquAvgLen.Enabled,
		len       : fc.EquAvgLen.CurValue,
		lenFo     : &fc.EquAvgLen,
	}

	if p.lenEnabled {
		p.len = p.lenFo.RandomValue()
	}

	return p
}

//=============================================================================

func (p *EquityVsAvgPart) Mutate() {

}

//=============================================================================

func (p *EquityVsAvgPart) Apply(f *db.TradingFilter) {

}

//=============================================================================
//===
//=== WinningPercPart
//===
//=============================================================================

type WinningPercPart struct {
	lenEnabled  bool
	len         int
	percEnabled bool
	perc        int

	lenFo  *optimization.FieldOptimization
	percFo *optimization.FieldOptimization
}

//=============================================================================

func NewWinningPercPart(fc *optimization.FilterConfig) Part {
	p := &WinningPercPart{
		lenEnabled : fc.WinPercLen.Enabled,
		len        : fc.WinPercLen.CurValue,
		lenFo      : &fc.WinPercLen,

		percEnabled: fc.WinPercPerc.Enabled,
		perc       : fc.WinPercPerc.CurValue,
		percFo     : &fc.WinPercPerc,
	}

	if p.lenEnabled {
		p.len = p.lenFo.RandomValue()
	}

	if p.percEnabled {
		p.perc = p.percFo.RandomValue()
	}

	return p
}

//=============================================================================

func (p *WinningPercPart) Mutate() {

}

//=============================================================================

func (p *WinningPercPart) Apply(f *db.TradingFilter) {

}

//=============================================================================
//===
//=== OldVsNewPart
//===
//=============================================================================

type OldVsNewPart struct {
	oldLenEnabled  bool
	oldLen         int

	newLenEnabled  bool
	newLen         int

	oldPercEnabled bool
	oldPerc        int

	oldLenFo  *optimization.FieldOptimization
	newLenFo  *optimization.FieldOptimization
	oldPercFo *optimization.FieldOptimization
}

//=============================================================================

func NewOldVsNewPart(fc *optimization.FilterConfig) Part {
	p := &OldVsNewPart{
		oldLenEnabled : fc.OldNewOldLen.Enabled,
		oldLen        : fc.OldNewOldLen.CurValue,
		oldLenFo      : &fc.OldNewOldLen,

		newLenEnabled : fc.OldNewNewLen.Enabled,
		newLen        : fc.OldNewNewLen.CurValue,
		newLenFo      : &fc.OldNewNewLen,

		oldPercEnabled: fc.OldNewOldPerc.Enabled,
		oldPerc       : fc.OldNewOldPerc.CurValue,
		oldPercFo     : &fc.OldNewOldPerc,
	}

	if p.oldLenEnabled {
		p.oldLen = p.oldLenFo.RandomValue()
	}

	if p.newLenEnabled {
		p.newLen = p.newLenFo.RandomValue()
	}

	if p.oldPercEnabled {
		p.oldPerc = p.oldPercFo.RandomValue()
	}

	return p
}

//=============================================================================

func (p *OldVsNewPart) Mutate() {

}

//=============================================================================

func (p *OldVsNewPart) Apply(f *db.TradingFilter) {

}


//=============================================================================
//===
//=== TrendlinePart
//===
//=============================================================================

type TrendlinePart struct {
	lenEnabled   bool
	len          int
	valueEnabled bool
	value        int

	lenFo   *optimization.FieldOptimization
	valueFo *optimization.FieldOptimization
}

//=============================================================================

func NewTrendlinePart(fc *optimization.FilterConfig) Part {
	p := &TrendlinePart{
		lenEnabled  : fc.TrendlineLen.Enabled,
		len         : fc.TrendlineLen.CurValue,
		lenFo       : &fc.TrendlineLen,

		valueEnabled: fc.TrendlineValue.Enabled,
		value       : fc.TrendlineValue.CurValue,
		valueFo     : &fc.TrendlineValue,
	}

	if p.lenEnabled {
		p.len = p.lenFo.RandomValue()
	}

	if p.valueEnabled {
		p.value = p.valueFo.RandomValue()
	}

	return p
}

//=============================================================================

func (p *TrendlinePart) Mutate() {

}

//=============================================================================

func (p *TrendlinePart) Apply(f *db.TradingFilter) {

}

//=============================================================================
