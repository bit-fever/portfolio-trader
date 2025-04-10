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

package business

import (
	"github.com/bit-fever/core/auth"
	"github.com/bit-fever/core/req"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func GetTradingSystems(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int, details bool) (*[]db.TradingSystem, error) {
	if ! c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	return db.GetTradingSystems(tx, filter, offset, limit)
}

//=============================================================================

func GetTrades(tx *gorm.DB, c *auth.Context, tsId uint) (*[]db.Trade, error) {
	ts, err := db.GetTradingSystemById(tx, tsId)
	if err != nil {
		return nil, err
	}
	if ts == nil {
		return nil, req.NewNotFoundError("Trading system not found: %v", tsId)
	}

	if ! c.Session.IsAdmin() {
		if ts.Username != c.Session.Username {
			return nil, req.NewForbiddenError("Trading system not owned by user: %v", tsId)
		}
	}

	return db.FindTradesByTsId(tx, tsId)
}

//=============================================================================

func DeleteTradingSystem(tx *gorm.DB, id uint) error {

	err := db.DeleteAllTradesByTradingSystemId(tx, id)
	if err != nil {
		return err
	}

	err = db.DeleteTradingFilter(tx, id)
	if err != nil {
		return err
	}

	return db.DeleteTradingSystem(tx, id)
}

//=============================================================================
