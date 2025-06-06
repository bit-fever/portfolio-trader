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
	"time"
)

//=============================================================================

func GetTradingSystems(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]TradingSystem, error) {
	var list []TradingSystem
	res := tx.Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetTradingSystemById(tx *gorm.DB, id uint) (*TradingSystem, error) {
	var list []TradingSystem
	res := tx.Find(&list, id)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	if len(list) == 1 {
		return &list[0], nil
	}

	return nil, nil
}

//=============================================================================

func GetTradingSystemsByUser(tx *gorm.DB, name string) (*[]TradingSystem, error) {
	var list []TradingSystem
	res := tx.Find(&list, "username = ?", name)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetUsersWithTradingSystems(tx *gorm.DB) ([]string, error) {
	var list []string
	res := tx.Table("trading_system").Distinct("username").Scan(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return list, nil
}

//=============================================================================

func GetTradingSystemsByIdsAsMap(tx *gorm.DB, ids []uint) (map[uint]*TradingSystem, error) {
	var list []TradingSystem
	res := tx.Find(&list, "id in ?", ids)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	tsMap := map[uint]*TradingSystem{}

	for _, ts := range list {
		tsAux := ts
		tsMap[ts.Id] = &tsAux
	}

	return tsMap, nil
}

//=============================================================================

func GetTradingSystemsInIdle(tx *gorm.DB, days int) (*[]TradingSystem, error) {
	var list []TradingSystem
	date := time.Now().Add(-time.Hour * 24 * time.Duration(days))

	res := tx.
		Where("running    = ?", true).
		Where("active     = ?", true).
		Where("last_trade < ?", date).
		Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func UpdateTradingSystem(tx *gorm.DB, ts *TradingSystem) error {
	return tx.Save(ts).Error
}

//=============================================================================

func UpdateDataProductInfo(tx *gorm.DB, dataProductId uint, values map[string]interface{}) error {
	return tx.Model(&TradingSystem{}).
		Where("data_product_id", dataProductId).
		Save(values).Error
}

//=============================================================================

func UpdateBrokerProductInfo(tx *gorm.DB, brokerProductId uint, values map[string]interface{}) error {
	return tx.Model(&TradingSystem{}).
		Where("broker_product_id", brokerProductId).
		Save(values).Error
}

//=============================================================================

func DeleteTradingSystem(tx *gorm.DB, id uint) error {
	return tx.Delete(&TradingSystem{}, id).Error
}

//=============================================================================
