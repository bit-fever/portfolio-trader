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

package db

import (
	"time"

	"github.com/bit-fever/core/req"
	"gorm.io/gorm"
)

//=============================================================================

func FindTradesByTsId(tx *gorm.DB, tsId uint) (*[]Trade, error) {
	var list []Trade

	filter := map[string]any{}
	filter["trading_system_id"] = tsId

	res := tx.Where(filter).Order("entry_date,exit_date").Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================
// We MUST order by entry_date because we may have cases like:
// entry_date,       exit_date
// 2025-1-1 03:00    2025-1-2 06:00
// 2025-1-2 06:00    2025-1-2 06:00    <-- fake trade
// 2025-1-2 06:00    2025-1-4 02:00
// Ordering by exit_date, the second record could come first

func FindTradesByTsIdFromTime(tx *gorm.DB, tsId uint, fromTime *time.Time, toTime *time.Time) (*[]Trade, error) {
	to   := time.Now().UTC()
	from := to.Add(-50 * 365 * 24 * time.Hour)

	if fromTime != nil {
		from = *fromTime
	}

	if toTime != nil {
		to = *toTime
	}

	var list []Trade

	//--- WHERE condition must be exit_date otherwise we loose trades started in the past and ended after fromTime
	query := "trading_system_id = ? and exit_date >= ? and exit_date<= ?"
	res   := tx.Order("exit_date,entry_date").Find(&list, query, tsId, from, to)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func FindTradesFromTime(tx *gorm.DB, tsIds []uint, fromTime time.Time) (*[]Trade, error) {
	var list []Trade

	res := tx.Find(&list, "trading_system_id in ? and entry_date >= ?", tsIds, fromTime)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func AddTrade(tx *gorm.DB, tr *Trade) error {
	err := tx.Create(tr).Error
	return req.NewServerErrorByError(err)
}

//=============================================================================

func DeleteAllTradesByTradingSystemId(tx *gorm.DB, id uint) error {
	err := tx.Delete(&Trade{}, "trading_system_id", id).Error
	return req.NewServerErrorByError(err)
}

//=============================================================================
