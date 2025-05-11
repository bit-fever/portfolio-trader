//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

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

import (
	"github.com/bit-fever/core/auth"
	"github.com/bit-fever/portfolio-trader/pkg/business/performance"
	"github.com/bit-fever/portfolio-trader/pkg/core"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
	"time"
)

//=============================================================================

func RunPerformanceAnalysis(tx *gorm.DB, c *auth.Context, tsId uint, req *performance.AnalysisRequest) (*performance.AnalysisResponse, error) {
	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	daysBack := req.DaysBack
	if daysBack == 0 {
		daysBack = 100000
	}
	t := time.Now()
	back := time.Hour * time.Duration(24 * daysBack)
	t = t.Add(-back)

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, t)
	if err != nil {
		return nil,err
	}

	res := performance.GetPerformanceAnalysis(ts, trades)

	//--- Get location and shift timezone according to user's preference

	loc, err := core.GetLocation(req.Timezone, res.TradingSystem)
	if err != nil {
		c.Log.Error("RunPerformanceAnalysis: Bad timezone", "timezone", req.Timezone, "error", err)
		return nil, err
	}

	res.AllEquities  .Time = shiftEquityTimezone(res.AllEquities  .Time, loc)
	res.LongEquities .Time = shiftEquityTimezone(res.LongEquities .Time, loc)
	res.ShortEquities.Time = shiftEquityTimezone(res.ShortEquities.Time, loc)
	res.Trades             = shiftTradesTimezone(res, loc)

	return res, nil
}

//=============================================================================

func shiftEquityTimezone(values *[]time.Time, loc *time.Location) *[]time.Time {

	var list []time.Time

	for _, tim := range *values {
		list = append(list, *shiftLocation(&tim, loc))
	}

	return &list
}

//=============================================================================

func shiftTradesTimezone(res *performance.AnalysisResponse, loc *time.Location) *[]db.Trade {

	var list []db.Trade

	for _, tr := range *res.Trades {
		tr.EntryDate         = shiftLocation(tr.EntryDate, loc)
		tr.ExitDate          = shiftLocation(tr.ExitDate, loc)
		tr.EntryDateAtBroker = shiftLocation(tr.EntryDateAtBroker, loc)
		tr.ExitDateAtBroker  = shiftLocation(tr.ExitDateAtBroker, loc)

		list = append(list, tr)
	}

	return &list
}

//=============================================================================

func shiftLocation(t *time.Time, loc *time.Location) *time.Time {
	if t == nil {
		return nil
	}

	out := t.In(loc)

	return &out
}

//=============================================================================
