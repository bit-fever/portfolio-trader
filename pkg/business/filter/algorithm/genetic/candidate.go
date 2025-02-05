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

import "github.com/bit-fever/portfolio-trader/pkg/business/filter/algorithm/optimization"

//=============================================================================

type Candidate struct {
	parts []Part
}

//=============================================================================

func NewRandomCandidate(fc *optimization.FilterConfig) *Candidate {
	c := &Candidate{
		parts: []Part{},
	}

	if fc.EnablePosProfit {
		c.parts = append(c.parts, NewPosProfitPart(fc))
	}

	if fc.EnableEquAvg {
		c.parts = append(c.parts, NewEquityVsAvgPart(fc))
	}

	if fc.EnableOldNew {
		c.parts = append(c.parts, NewOldVsNewPart(fc))
	}

	if fc.EnableWinPerc {
		c.parts = append(c.parts, NewWinningPercPart(fc))
	}

	if fc.EnableTrendline {
		c.parts = append(c.parts, NewTrendlinePart(fc))
	}

	return c
}

//=============================================================================

func (c *Candidate) CrossOver(c2 *Candidate) *Candidate {
	//TODO
	return nil
}

//=============================================================================

func (c *Candidate) Mutate() {

}

//=============================================================================
