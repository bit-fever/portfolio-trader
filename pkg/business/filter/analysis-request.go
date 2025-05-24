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

import "time"

//=============================================================================
//===
//=== AnalysisRequest
//===
//=============================================================================

type AnalysisRequest struct {
	StartDate *time.Time      `json:"startDate,omitempty"`
	Filter    *TradingFilter  `json:"filter,omitempty"`
}

//=============================================================================

type TradingFilter struct {
	EquAvgEnabled    bool   `json:"equAvgEnabled"`
	EquAvgLen        int    `json:"equAvgLen"`
	PosProEnabled    bool   `json:"posProEnabled"`
	PosProLen        int    `json:"posProLen"`
	WinPerEnabled    bool   `json:"winPerEnabled"`
	WinPerLen        int    `json:"winPerLen"`
	WinPerValue      int    `json:"winPerValue"`
	OldNewEnabled    bool   `json:"oldNewEnabled"`
	OldNewOldLen     int    `json:"oldNewOldLen"`
	OldNewOldPerc    int    `json:"oldNewOldPerc"`
	OldNewNewLen     int    `json:"oldNewNewLen"`
	TrendlineEnabled bool   `json:"trendlineEnabled"`
	TrendlineLen     int    `json:"trendlineLen"`
	TrendlineValue   int    `json:"trendlineValue"`
	DrawdownEnabled  bool   `json:"drawdownEnabled"`
	DrawdownMin      int    `json:"drawdownMin"`
	DrawdownMax      int    `json:"drawdownMax"`
}

//=============================================================================
