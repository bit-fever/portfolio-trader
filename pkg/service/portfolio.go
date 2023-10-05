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
	"github.com/bit-fever/portfolio-trader/pkg/business"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"github.com/bit-fever/portfolio-trader/pkg/tool"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//=============================================================================

func getPortfolios(c *gin.Context) {
	offset, limit, err := tool.GetPagingParams(c)

	if err == nil {
		err = db.RunInTransaction(func(tx *gorm.DB) error {
			list, err := db.GetPortfolios(tx, nil, offset, limit)

			if err != nil {
				return err
			}

			return tool.ReturnList(c, list, offset, limit, len(*list))
		})
	}

	tool.ReturnError(c, err)
}

//=============================================================================

func getPortfolioTree(c *gin.Context) {
	err := db.RunInTransaction(func(tx *gorm.DB) error {
		list, err := business.GetPortfolioTree(tx)

		if err != nil {
			return err
		}

		return tool.ReturnObject(c, list)
	})

	tool.ReturnError(c, err)
}

//=============================================================================

func getPortfolioMonitoring(c *gin.Context) {
	params := business.PortfolioMonitoringParams{}
	err    := tool.BindParamsFromBody(c, &params)

	if err == nil {
		err = db.RunInTransaction(func(tx *gorm.DB) error {
			result, err := business.GetPortfolioMonitoring(tx, &params)

			if err != nil {
				return err
			}

			return tool.ReturnObject(c, result)
		})
	}

	tool.ReturnError(c, err)
}

//=============================================================================
