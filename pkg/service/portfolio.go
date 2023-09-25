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
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"github.com/bit-fever/portfolio-trader/pkg/model"
	"github.com/bit-fever/portfolio-trader/pkg/tool"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

//=============================================================================

func getPortfolios(c *gin.Context) {
	offset, limit, err := tool.GetPagingParams(c)

	if err != nil {
		return
	}

	err = db.RunInTransaction(func(tx *gorm.DB) error {
		list, err := db.GetPortfolios(tx, nil, offset, limit)

		if err != nil {
			return err
		}

		return tool.Return(c, list, offset, limit, len(*list))
	})

	if err != nil {
		tool.ErrorInternal(c, err.Error())
	}
}

//=============================================================================

func getPortfolioTree(c *gin.Context) {
	err := db.RunInTransaction(func(tx *gorm.DB) error {
		list, err := db.GetPortfolioTree(tx)

		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, list)
		return nil
	})

	if err != nil {
		tool.ErrorInternal(c, err.Error())
	}
}

//=============================================================================

func getPortfolioMonitoring(c *gin.Context) {
	params := model.PortfolioMonitoringParams{}
	if err := tool.BindParamsFromBody(c, &params); err != nil {
		return
	}

	err := db.RunInTransaction(func(tx *gorm.DB) error {
		result, err := db.GetPortfolioMonitoring(tx, &params)

		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, result)
		return nil
	})

	if err != nil {
		_ = tool.ErrorInternal(c, err.Error())
	}
}

//=============================================================================
