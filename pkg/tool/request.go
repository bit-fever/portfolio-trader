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

package tool

import (
	"github.com/bit-fever/portfolio-trader/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//=============================================================================

func Error(c *gin.Context, errorCode int, errorMessage string, details any) {
	c.JSON(errorCode, &model.ErrorResponse{
		Code:    errorCode,
		Error:   errorMessage,
		Details: details,
	})
}

//=============================================================================

func ErrorBadRequest(c *gin.Context, errorMessage string, details any) {
	Error(c, http.StatusBadRequest, errorMessage, details)
}

//=============================================================================

func ErrorInternal(c *gin.Context, errorMessage string) {
	Error(c, http.StatusInternalServerError, errorMessage, nil)
}

//=============================================================================

func GetPagingParams(c *gin.Context) (offset int, limit int, errV error) {

	params := c.Request.URL.Query()

	//--- Extract offset

	if !params.Has("offset") {
		offset = 0
	} else {
		offsetP := params.Get("offset")
		offsetV, err := strconv.ParseInt(offsetP, 10, 64)

		if err != nil || offsetV < 0 {
			ErrorBadRequest(c, "Invalid offset param", offsetP)
			return 0, 0, err
		}

		offset = int(offsetV)
	}

	//--- Extract limit

	if !params.Has("limit") {
		limit = model.MaxLimit
	} else {
		limitP := params.Get("limit")
		limitV, err := strconv.ParseInt(limitP, 10, 32)

		if err != nil || limitV < 1 || limit > model.MaxLimit {
			ErrorBadRequest(c, "Invalid limit param", limitP)
			return 0, 0, err
		}

		limit = int(limitV)
	}

	return offset, limit, nil
}

//=============================================================================

func Return(c *gin.Context, result any, offset int, limit int, size int) error {
	c.JSON(http.StatusOK, &model.ListResponse{
		Offset:   offset,
		Limit:    limit,
		Overflow: size == model.MaxLimit,
		Result:   result,
	})

	return nil
}

//=============================================================================
