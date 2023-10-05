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

func GetTradingSystems(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]TradingSystem, error) {
	var list []TradingSystem
	res := tx.Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, tool.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetTradingSystemById(tx *gorm.DB, id uint) (*TradingSystem, error) {
	var list []TradingSystem
	res := tx.Find(&list, id)

	if res.Error != nil {
		return nil, tool.NewServerErrorByError(res.Error)
	}

	if len(list) == 1 {
		return &list[0], nil
	}

	return nil, nil
}

//=============================================================================

func GetTradingSystemsById(tx *gorm.DB, ids []uint) (*[]TradingSystem, error) {
	var list []TradingSystem
	res := tx.Find(&list, ids)

	if res.Error != nil {
		return nil, tool.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetTradingSystemsByIdAsMap(tx *gorm.DB, ids []uint) (map[int]*TradingSystem, error) {
	var list []TradingSystem
	res := tx.Find(&list, ids)

	if res.Error != nil {
		return nil, tool.NewServerErrorByError(res.Error)
	}

	tsMap := map[int]*TradingSystem{}

	for _, ts := range list {
		tsAux := ts
		tsMap[int(ts.Id)] = &tsAux
	}

	return tsMap, nil
}

//=============================================================================

func GetTradingSystemsFull(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]TradingSystemFull, error) {
	var list []TradingSystemFull
	query := "SELECT ts.*, i.ticker as instrument_ticker, p.name as portfolio_name " +
		"FROM trading_system ts " +
		"LEFT JOIN instrument i on ts.instrument_id = i.id " +
		"LEFT JOIN portfolio p on ts.portfolio_id = p.id"

	res := tx.Raw(query).Find(&list)

	if res.Error != nil {
		return nil, tool.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetOrCreateTradingSystem(tx *gorm.DB, code string, ts *TradingSystem) (*TradingSystem, error) {
	res := tx.Where(&TradingSystem{Code: code}).FirstOrCreate(&ts)

	if res.Error != nil {
		return nil, tool.NewServerErrorByError(res.Error)
	}

	return ts, nil
}

//=============================================================================
