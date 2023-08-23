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
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

//=============================================================================

func getTradingSystems(c *gin.Context) {
	err := db.RunInTransaction(func(tx *gorm.DB) error {
		list, err := db.GetTradingSystems(tx)

		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, &list)
		return nil
	})

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error)
	}
}

//=============================================================================

func getDailyInfo(c *gin.Context) {
	tsId, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error)
	}

	err = db.RunInTransaction(func(tx *gorm.DB) error {
		list, err := db.GetDailyInfo(tx, int(tsId))

		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, &list)
		return nil
	})

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error)
	}
}

//=============================================================================
