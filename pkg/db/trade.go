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
	"github.com/bit-fever/core/req"
	"gorm.io/gorm"
	"time"
)

//=============================================================================

func FindTradesByTsId(tx *gorm.DB, tsId uint) (*[]Trade, error) {
	var list []Trade

	filter := map[string]any{}
	filter["trading_system_id"] = tsId

	res := tx.Where(filter).Find(&list).Order("entry_date")

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func FindTradesFromTime(tx *gorm.DB, tsIds []uint, fromTime time.Time) (*[]Trade, error) {
	var data []Trade

	res := tx.Find(&data, "trading_system_id in ? and exit_time >= ?", tsIds, fromTime)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &data, nil
}

//=============================================================================

func FindLastTrade(tx *gorm.DB, tsId uint) (*Trade, error) {
	var data []Trade

	res := tx.Order("entry_time DESC").Limit(1).Find(&data, "trading_system_id = ?", tsId)
	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	if len(data) == 0 {
		return nil, nil
	}

	return &data[0], nil
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
