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

package business

import (
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"github.com/bit-fever/portfolio-trader/pkg/tool"
	"gorm.io/gorm"
)

//=============================================================================

func GetFilteringAnalysis(tx *gorm.DB, tsId uint, params *FilteringParams) (*FilteringResponse, error) {

	ts, err := db.GetTradingSystemById(tx, tsId)
	if err != nil {
		return nil, err
	}

	if ts == nil {
		return nil, tool.NewRequestError("Missing trading system with id=%v", tsId)
	}

	if params.IsEmpty() {
		//if err := json.Unmarshal([]byte(ts.Filter), &params); err != nil {
		//	return nil, errors.New("Bad json format for filter : %v", ts.Filter)
		//}
	}

	res := &FilteringResponse{}

	return res, nil
}

//=============================================================================
