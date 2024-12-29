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

import "github.com/bit-fever/portfolio-trader/pkg/db"

//=============================================================================

type Part interface {
	Mutate()
	Apply(filter *db.TradingFilter)
}

//=============================================================================

type PartConfig struct {
	minValue int
	maxValue int
	step     int
}

//=============================================================================
//===
//=== PosProfitPart
//===
//=============================================================================

type PosProfitPart struct {
	enabled  bool
	len      int
	config   PartConfig
}

//=============================================================================

func NewPosProfitPart() Part {
	return nil
}

//=============================================================================
//===
//=== EquityVsAvgPart
//===
//=============================================================================

type EquityVsAvgPart struct {

}

//=============================================================================

func NewEquityVsAvgPart() Part {
	return nil
}

//=============================================================================
//===
//=== WinningPercPart
//===
//=============================================================================

type WinningPercPart struct {

}

//=============================================================================

func NewWinningPercPart() Part {
	return nil
}

//=============================================================================
//===
//=== OldVsNewPart
//===
//=============================================================================

type OldVsNewPart struct {

}

//=============================================================================

func NewOldVsNewPart() Part {
	return nil
}

//=============================================================================
//===
//=== TrendlinePart
//===
//=============================================================================

type TrendlinePart struct {

}

//=============================================================================

func NewTrendlinePart() Part {
	return nil
}

//=============================================================================
