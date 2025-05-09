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
	"github.com/bit-fever/core/auth"
	"github.com/bit-fever/portfolio-trader/pkg/business"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func getPortfolios(c *auth.Context) {
	filter := map[string]any{}
	offset, limit, err := c.GetPagingParams()

	if err == nil {
		err = db.RunInTransaction(func(tx *gorm.DB) error {
			list, err := business.GetPortfolios(tx, c, filter, offset, limit)

			if err != nil {
				return err
			}

			return c.ReturnList(list, offset, limit, len(*list))
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func getPortfolioTree(c *auth.Context) {
	filter := map[string]any{}
	offset, limit, err := c.GetPagingParams()

	if err == nil {
		err = db.RunInTransaction(func(tx *gorm.DB) error {
			list, err := business.GetPortfolioTree(tx, c, filter, offset, limit)

			if err != nil {
				return err
			}

			return c.ReturnObject(list)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func getPortfolioMonitoring(c *auth.Context) {
	params := business.PortfolioMonitoringParams{}
	err    := c.BindParamsFromBody(&params)

	if err == nil {
		err = db.RunInTransaction(func(tx *gorm.DB) error {
			result, err := business.GetPortfolioMonitoring(tx, &params)

			if err != nil {
				return err
			}

			return c.ReturnObject(result)
		})
	}

	c.ReturnError(err)
}

//=============================================================================
