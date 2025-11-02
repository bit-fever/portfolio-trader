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

package service

import (
	"github.com/tradalia/core/auth"
	"github.com/tradalia/portfolio-trader/pkg/business"
	"github.com/tradalia/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func setTradingSystemTrading(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		req := business.TradingSystemTradingRequest{}
		err = c.BindParamsFromBody(&req)

		if err == nil {
			err = db.RunInTransaction(func(tx *gorm.DB) error {
				rep, err := business.SetTradingSystemTrading(tx, c, tsId, &req)

				if err != nil {
					return err
				}

				return c.ReturnObject(rep)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func setTradingSystemRunning(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		req := business.TradingSystemRunningRequest{}
		err = c.BindParamsFromBody(&req)

		if err == nil {
			err = db.RunInTransaction(func(tx *gorm.DB) error {
				rep, err := business.SetTradingSystemRunning(tx, c, tsId, &req)

				if err != nil {
					return err
				}

				return c.ReturnObject(rep)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func setTradingSystemActivation(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		req := business.TradingSystemActivationRequest{}
		err = c.BindParamsFromBody(&req)

		if err == nil {
			err = db.RunInTransaction(func(tx *gorm.DB) error {
				rep, err := business.SetTradingSystemActivation(tx, c, tsId, &req)

				if err != nil {
					return err
				}

				return c.ReturnObject(rep)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func setTradingSystemActive(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		req := business.TradingSystemActiveRequest{}
		err = c.BindParamsFromBody(&req)

		if err == nil {
			err = db.RunInTransaction(func(tx *gorm.DB) error {
				rep, err := business.SetTradingSystemActive(tx, c, tsId, &req)

				if err != nil {
					return err
				}

				return c.ReturnObject(rep)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================
