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
	"time"

	"github.com/bit-fever/core/auth"
	"github.com/bit-fever/core/datatype"
	"github.com/bit-fever/core/req"
	"github.com/bit-fever/portfolio-trader/pkg/business/performance"
	"github.com/bit-fever/portfolio-trader/pkg/core"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func RunPerformanceAnalysis(tx *gorm.DB, c *auth.Context, tsId uint, req *performance.AnalysisRequest) (*performance.AnalysisResponse, error) {

	//--- Get trading system

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	//--- Get location of timezone to shift dates

	loc, err := core.GetLocation(req.Timezone, ts)
	if err != nil {
		c.Log.Error("RunPerformanceAnalysis: Bad timezone", "timezone", req.Timezone, "error", err)
		return nil, err
	}

	fromTime, toTime, err := calcPerformancePeriod(req.DaysBack, req.FromDate, req.ToDate, loc)
	if err != nil {
		c.Log.Error("RunPerformanceAnalysis: Bad fromDate or toDate", "fromDate", req.FromDate, "toDate", req.ToDate, "error", err)
		return nil, err
	}

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, fromTime, toTime)
	if err != nil {
		return nil,err
	}
	shiftTradesTimezone(trades, loc)

	returns,err := db.FindDailyReturnsByTsIdFromTime(tx, ts.Id, fromTime, toTime)
	if err != nil {
		return nil,err
	}

	res := performance.GetPerformanceAnalysis(ts, trades, returns)

	return res, nil
}

//=============================================================================

func calcPerformancePeriod(daysBack int, fromDate, toDate datatype.IntDate, loc *time.Location) (*time.Time, *time.Time, error) {
	//--- All

	if daysBack == 0 {
		return nil, nil, nil
	}

	//--- Specific last days

	if daysBack > 0 {
		fromTime := time.Now().UTC()
		back     := time.Hour * time.Duration(24 * daysBack)
		fromTime = fromTime.Add(-back)

		return &fromTime, nil, nil
	}

	//--- Custom range

	if daysBack == -1 {
		var from *time.Time
		var to   *time.Time

		if !fromDate.IsNil() {
			if !fromDate.IsValid() {
				return nil, nil, req.NewBadRequestError("Invalid fromDate parameter: %d", fromDate)
			}

			tt := fromDate.ToDateTime(false, loc)
			from = &tt
		}

		if !toDate.IsNil() {
			if !toDate.IsValid() {
				return nil, nil, req.NewBadRequestError("Invalid toDate parameter: %d", toDate)
			}

			tt := toDate.ToDateTime(true, loc)
			to = &tt
		}

		return from, to, nil
	}

	return nil, nil, req.NewBadRequestError("Invalid daysBack parameter: %d", daysBack)
}

//=============================================================================

func shiftTradesTimezone(trades *[]db.Trade, loc *time.Location) {
	for i:=0; i<len(*trades); i++ {
		tr := &(*trades)[i]
		tr.EntryDate         = shiftLocation(tr.EntryDate,         loc)
		tr.ExitDate          = shiftLocation(tr.ExitDate,          loc)
		tr.EntryDateAtBroker = shiftLocation(tr.EntryDateAtBroker, loc)
		tr.ExitDateAtBroker  = shiftLocation(tr.ExitDateAtBroker,  loc)
	}
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
