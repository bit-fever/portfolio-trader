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
	"github.com/bit-fever/portfolio-trader/pkg/tool"
	"gorm.io/gorm"
)

//=============================================================================

func FindDailyInfoByTsId(tx *gorm.DB, tsId uint) ([]TsDailyInfo, error) {
	var list []TsDailyInfo

	filter := map[string]any{}
	filter["trading_system_id"] = tsId

	res := tx.Where(filter).Find(&list)

	if res.Error != nil {
		return nil, tool.NewServerErrorByError(res.Error)
	}

	return list, nil
}

//=============================================================================

func FindDailyInfoByTsIdAsMap(tx *gorm.DB, tsId uint) (map[int]TsDailyInfo, error) {
	list,err := FindDailyInfoByTsId(tx, tsId)

	if err != nil {
		return nil, err
	}

	diMap := map[int]TsDailyInfo{}

	for _, di := range list {
		diMap[di.Day] = di
	}

	return diMap, nil
}

//=============================================================================

func FindDailyInfoFromDay(tx *gorm.DB, tsIds []uint, fromDay int) (*[]TsDailyInfo, error) {
	var data []TsDailyInfo

	res := tx.Find(&data, "trading_system_id in ? and day >= ?", tsIds, fromDay)

	if res.Error != nil {
		return nil, tool.NewServerErrorByError(res.Error)
	}

	return &data, nil
}

//=============================================================================

func AddTsDailyInfo(tx *gorm.DB, di *TsDailyInfo) error {
	err := tx.Create(di).Error
	return tool.NewServerErrorByError(err)
}

//=============================================================================
