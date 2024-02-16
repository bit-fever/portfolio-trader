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

//=============================================================================
//===
//=== Activation
//===
//=============================================================================

type Activation struct {
	Days   []int   `json:"days"`
	Values []int8  `json:"values"`
}

//-----------------------------------------------------------------------------

func (p *Activation) AddDay(day int, value int8) {
	p.Days   = append(p.Days,   day)
	p.Values = append(p.Values, value)
}

//=============================================================================
//===
//=== ActivationStrategy
//===
//=============================================================================

type ActivationStrategy struct {
	activation *Activation
	enabled    bool
	index      int
}

//=============================================================================

func (as *ActivationStrategy) IsActive(day int) bool {
	//--- Strategy not enabled: skip it returning always 1
	if !as.enabled {
		return true
	}

	//--- Strategy not computable: return true because we must align with unfiltered equity
	if as.activation == nil {
		return true
	}

	if day<as.activation.Days[as.index] {
		return true
	}

	if day != as.activation.Days[as.index] {
		panic("Help!")
	}

	as.index++
	return as.activation.Values[as.index -1] != 0
}

//=============================================================================

func NewActivationStrategy(a *Activation, enabled bool) *ActivationStrategy {
	return &ActivationStrategy{
		activation: a,
		enabled: enabled,
		index: 0,
	}
}

//=============================================================================
