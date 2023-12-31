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

package business

import "github.com/bit-fever/core/req"

//=============================================================================

type EquityFilter interface {
	init(c *FilteringConfig) error
	compute(e *Equities, index int) bool
}

//=============================================================================

func createFilterChain(c *FilteringConfig) (*[]EquityFilter, error) {
	var list []EquityFilter

	//--- Append equity average filter

	if c.EquityAverage.Enabled {
		eaf := EquityAverageFilter{}
		if err:=eaf.init(c); err != nil {
			return nil, err
		}

		list = append(list, &eaf)
	}

	//--- Append long/short period filter

	if c.LongShort.Enabled {
		lsf := LongShortPeriodFilter{}
		if err:=lsf.init(c); err != nil {
			return nil, err
		}

		list = append(list, &lsf)
	}

	return &list, nil
}

//=============================================================================
//===
//=== Equity average filtering
//===
//=============================================================================

type EquityAverageFilter struct {
	days int
}

//=============================================================================

func (e *EquityAverageFilter) init(c *FilteringConfig) error {
	e.days = c.EquityAverage.Days

	if e.days <1 || e.days > 200 {
		return req.NewBadRequestError("Invalid range for Average days (must be 1..200): %d", e.days)
	}

	return nil
}

//=============================================================================

func (e *EquityAverageFilter) compute(eq *Equities, index int) bool {
	return eq.UnfilteredProfit[index] >= eq.Average[index]
}

//=============================================================================
//===
//=== Long/Short periods filtering
//===
//=============================================================================

type LongShortPeriodFilter struct {}

//=============================================================================

func (e *LongShortPeriodFilter) init(c *FilteringConfig) error {
	return nil
}

//=============================================================================

func (e *LongShortPeriodFilter) compute(eq *Equities, index int) bool {
	return true
}

//=============================================================================
