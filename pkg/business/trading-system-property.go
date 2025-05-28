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
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

type TradingSystemTradingRequest struct {
	Value bool `json:"value"`
}

//=============================================================================

type TradingSystemRunningRequest struct {
	Value bool `json:"value"`
}

//=============================================================================

type TradingSystemActivationRequest struct {
	Value bool `json:"value"`
}

//=============================================================================

type TradingSystemActiveRequest struct {
	Value bool `json:"value"`
}

//=============================================================================

const (
	ResponseStatusOk     = "ok"
	ResponseStatusSkipped= "skipped"
	ResponseStatusError  = "error"
)

//-----------------------------------------------------------------------------

type TradingSystemPropertyResponse struct {
	Status        string            `json:"status"`
	Message       string            `json:"message"`
	TradingSystem *db.TradingSystem `json:"tradingSystem"`
}

//=============================================================================

func SetTradingSystemTrading(tx *gorm.DB, c *auth.Context, tsId uint, req *TradingSystemTradingRequest) (*TradingSystemPropertyResponse, error) {
	c.Log.Info("SetTradingSystemTrading: Trading property change request", "id", tsId, "value", req.Value)

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	oldValue := ts.Trading
	newValue := req.Value

	if oldValue == newValue {
		return &TradingSystemPropertyResponse{
			Status: ResponseStatusSkipped,
		}, nil
	}

	if !oldValue && newValue {
		//--- Turning on
	} else {
		//--- Turning off
		if ts.Running {
			return &TradingSystemPropertyResponse{
				Status : ResponseStatusError,
				Message: "Trading system must be stopped",
			}, nil
		}
	}

	ts.Trading = newValue
	updateStatus(ts)
	err = db.UpdateTradingSystem(tx, ts)
	if err != nil {
		return nil, err
	}

	c.Log.Info("SetTradingSystemTrading: Trading property changed", "id", tsId, "value", req.Value)

	return &TradingSystemPropertyResponse{
		Status       : ResponseStatusOk,
		TradingSystem: ts,
	}, err
}

//=============================================================================

func SetTradingSystemRunning(tx *gorm.DB, c *auth.Context, tsId uint, req *TradingSystemRunningRequest) (*TradingSystemPropertyResponse, error) {
	c.Log.Info("SetTradingSystemRunning: Running property change request", "id", tsId, "value", req.Value)

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	oldValue := ts.Running
	newValue := req.Value

	if oldValue == newValue {
		return &TradingSystemPropertyResponse{
			Status : ResponseStatusSkipped,
		}, nil
	}

	ts.Running = newValue
	updateStatus(ts)
	err = db.UpdateTradingSystem(tx, ts)
	if err != nil {
		return nil, err
	}

	err = updateRewind(ts)
	if err != nil {
		return nil, err
	}

	c.Log.Info("SetTradingSystemRunning: Running property changed", "id", tsId, "value", req.Value)

	return &TradingSystemPropertyResponse{
		Status       : ResponseStatusOk,
		TradingSystem: ts,
	}, err
}

//=============================================================================

func SetTradingSystemActivation(tx *gorm.DB, c *auth.Context, tsId uint, req *TradingSystemActivationRequest) (*TradingSystemPropertyResponse, error) {
	c.Log.Info("SetTradingSystemActivation: Auto-activation property change request", "id", tsId, "value", req.Value)

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	oldValue := ts.AutoActivation
	newValue := req.Value

	if oldValue == newValue {
		return &TradingSystemPropertyResponse{
			Status : ResponseStatusSkipped,
		}, nil
	}

	ts.AutoActivation = newValue
	err = db.UpdateTradingSystem(tx, ts)

	c.Log.Info("SetTradingSystemActivation: Auto-activation property changed", "id", tsId, "value", req.Value)

	return &TradingSystemPropertyResponse{
		Status       : ResponseStatusOk,
		TradingSystem: ts,
	}, err
}

//=============================================================================

func SetTradingSystemActive(tx *gorm.DB, c *auth.Context, tsId uint, req *TradingSystemActiveRequest) (*TradingSystemPropertyResponse, error) {
	c.Log.Info("SetTradingSystemActive: Active property change request", "id", tsId, "value", req.Value)

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	oldValue := ts.Active
	newValue := req.Value

	if oldValue == newValue {
		return &TradingSystemPropertyResponse{
			Status : ResponseStatusSkipped,
		}, nil
	}

	if ts.AutoActivation {
		return &TradingSystemPropertyResponse{
			Status : ResponseStatusError,
			Message: "Trading system is in AUTOMATIC mode. Switch to MANUAL to change",
		}, nil
	}

	ts.Active = newValue
	updateStatus(ts)
	err = db.UpdateTradingSystem(tx, ts)
	if err != nil {
		return nil, err
	}

	err = updateRewind(ts)

	c.Log.Info("SetTradingSystemActive: Active property changed", "id", tsId, "value", req.Value)

	return &TradingSystemPropertyResponse{
		Status       : ResponseStatusOk,
		TradingSystem: ts,
	}, err
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func updateStatus(ts *db.TradingSystem) {
	ts.SuggestedAction = db.TsActionNone

	if ! ts.Running {
		ts.Status = db.TsStatusOff
	} else if ts.Active {
		ts.Status = db.TsStatusRunning
	} else {
		ts.Status = db.TsStatusPaused
	}
}

//============================================================================

func updateRewind(ts *db.TradingSystem) error {
	return nil
}

//=============================================================================
