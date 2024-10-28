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

package db

import (
	"github.com/bit-fever/core/req"
	"gorm.io/gorm"
)

//=============================================================================

func GetTradingFiltersByTsId(tx *gorm.DB, tsId uint) (*TradingFilter, error) {
	var list []TradingFilter

	filter := map[string]any{}
	filter["trading_system_id"] = tsId

	res := tx.Where(filter).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	if len(list) == 0 {
		return &TradingFilter{
			TradingSystemId : tsId,
			EquAvgLen       :   30,
			PosProLen       :   45,
			WinPerLen       :   45,
			WinPerValue     :   30,
			OldNewOldLen    :   45,
			OldNewOldPerc   :   90,
			OldNewNewLen    :   45,
		}, nil
	}

	return &list[0], nil
}

//=============================================================================

func GetTradingFiltersByTsIds(tx *gorm.DB, tsIds []uint) (*[]TradingFilter, error) {
	var list []TradingFilter

	filter := map[string]any{}
	filter["trading_system_id"] = tsIds

	res := tx.Where(filter).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func SetTradingFilters(tx *gorm.DB, tf *TradingFilter) {
	tx.Save(tf)
}

//=============================================================================
