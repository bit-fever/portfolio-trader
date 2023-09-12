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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bit-fever/portfolio-trader/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//=============================================================================

func Error(c *gin.Context, errorCode int, errorMessage string, details any) error {
	c.JSON(errorCode, &model.ErrorResponse{
		Code:    errorCode,
		Error:   errorMessage,
		Details: details,
	})

	return errors.New(errorMessage)
}

//=============================================================================

func ErrorBadRequest(c *gin.Context, errorMessage string, details any) error {
	return Error(c, http.StatusBadRequest, errorMessage, details)
}

//=============================================================================

func ErrorInternal(c *gin.Context, errorMessage string) error {
	return Error(c, http.StatusInternalServerError, errorMessage, nil)
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

func BindParamsFromQuery(c *gin.Context, obj any) (err error) {
	if err := c.ShouldBindQuery(obj); err != nil {
		message := parseError(err)
		_ = ErrorBadRequest(c, message, nil)
		return errors.New(message)
	}

	return nil
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func parseError(err error) string {
	switch typedError := any(err).(type) {
	case validator.ValidationErrors:
		for _, e := range typedError {
			return parseFieldError(e)
		}

	case *json.UnmarshalTypeError:
		return parseMarshallingError(*typedError)

	case *strconv.NumError:
		return parseConvertError(*typedError)
	}

	return err.Error()
}

//=============================================================================

func parseFieldError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())
	fieldPrefix := fmt.Sprintf("The field %s", field)
	tag := strings.Split(e.Tag(), "|")[0]

	switch tag {
	case "required":
		return fmt.Sprintf("Missing the '%s' parameter", field)

	case "required_without":
		return fmt.Sprintf("%s is required if %s is not supplied", fieldPrefix, e.Param())

	case "lt", "ltfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be less than %s", fieldPrefix, param)

	case "gt", "gtfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be greater than %s", fieldPrefix, param)

	default:
		return fmt.Errorf("%v", e).Error()
	}
}

//=============================================================================

func parseMarshallingError(e json.UnmarshalTypeError) string {
	return fmt.Sprintf("Invalid type: '%s' must be a %s", strings.ToLower(e.Field), e.Type.String())
}

//=============================================================================

func parseConvertError(e strconv.NumError) string {
	return fmt.Sprintf("Parameter must be an integer: %s", e.Num)
}

//=============================================================================