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
	"github.com/bit-fever/core/auth"
	"github.com/bit-fever/core/req"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
	"log/slog"
)

//=============================================================================

func GetPortfolios(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int) (*[]db.Portfolio, error) {
	if ! c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	return db.GetPortfolios(tx, filter, offset, limit)
}

//=============================================================================

func GetPortfolioTree(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int) (*[]*PortfolioTree, error) {

	//--- The only valid filter can be the username

	//--- Get all portfolios

	poList, err := GetPortfolios(tx, c, filter, offset, limit)
	if err != nil {
		return nil, req.NewServerErrorByError(err)
	}

	//--- Get all trading systems

	tsList, err := GetTradingSystems(tx, c, filter, offset, limit, false)
	if err != nil {
		return nil, req.NewServerErrorByError(err)
	}

	return buildPortfolioTree(c.Log, poList, tsList), nil
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func buildPortfolioTree(log *slog.Logger, poList *[]db.Portfolio, tsList *[]db.TradingSystemFull) *[]*PortfolioTree {

	//--- Step 1: Collect all nodes into a map

	nodeMap := map[uint]*PortfolioTree{}
	fullMap := map[uint]*PortfolioTree{}

	for _, p := range *poList {
		pt := &PortfolioTree{
			Portfolio:      p,
			Children:       []*PortfolioTree{},
			TradingSystems: []*db.TradingSystemFull{},
		}
		nodeMap[p.Id] = pt
		fullMap[p.Id] = pt
	}

	//--- Step 2: Build the tree

	for key, p := range fullMap {
		if p.ParentId != 0 {
			parent := fullMap[p.ParentId]
			parent.AddChild(p)
			delete(nodeMap, key)
		}
	}

	//--- Step 2: Add trading system information

	for _, ts := range *tsList {
		aux := ts
		portfolio := fullMap[ts.PortfolioId]
		portfolio.AddTradingSystem(&aux)
	}

	//--- Step 3: Return tree

	if len(*poList) > 0 && len(nodeMap) == 0 {
		log.Error("Portfolios have circular loops (!)")
	}

	var result []*PortfolioTree

	for _, p := range nodeMap {
		result = append(result, p)
	}

	return &result
}

//=============================================================================
