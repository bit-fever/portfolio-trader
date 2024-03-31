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

//=============================================================================
//===
//=== AnalysisRequest
//===
//=============================================================================

type AnalysisRequest struct {
	Filters *TradingFilters  `json:"filters,omitempty"`
}

//=============================================================================

type TradingFilters struct {
	EquAvgEnabled   bool   `json:"equAvgEnabled"`
	EquAvgDays      int    `json:"equAvgDays"`
	PosProEnabled   bool   `json:"posProEnabled"`
	PosProDays      int    `json:"posProDays"`
	WinPerEnabled   bool   `json:"winPerEnabled"`
	WinPerDays      int    `json:"winPerDays"`
	WinPerValue     int    `json:"winPerValue"`
	OldNewEnabled   bool   `json:"oldNewEnabled"`
	OldNewOldDays   int    `json:"oldNewOldDays"`
	OldNewOldPerc   int    `json:"oldNewOldPerc"`
	OldNewNewDays   int    `json:"oldNewNewDays"`
}

//=============================================================================
