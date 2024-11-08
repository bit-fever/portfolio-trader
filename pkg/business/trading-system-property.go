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
	"errors"
	"github.com/bit-fever/core/auth"
	"github.com/bit-fever/portfolio-trader/pkg/core/tradingsystem"
	"github.com/bit-fever/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

const (
	TsPropertyRunning    = "running"
	TsPropertyActivation = "activation"
	TsPropertyActive     = "active"
)

//-----------------------------------------------------------------------------

type TradingSystemPropertyRequest struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

//=============================================================================

const (
	ResponseStatusOk     = "ok"
	ResponseStatusSkipped= "skipped"
	ResponseStatusError  = "error"
)

//-----------------------------------------------------------------------------

type TradingSystemPropertyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

//=============================================================================

func SetTradingSystemProperty(tx *gorm.DB, c *auth.Context, tsId uint, req *TradingSystemPropertyRequest) (*TradingSystemPropertyResponse, error) {
	c.Log.Info("SetTradingSystemProperty: Property change request", "tsId", tsId, "property", req.Property, "value", req.Value)

	ts, err := getTradingSystem(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	switch req.Property {
		case TsPropertyRunning:
			return handleRunningProperty(tx, c, ts, req.Value)
		case TsPropertyActivation:
			return handleActivationProperty(tx, c, ts, req.Value)
		case TsPropertyActive:
			return handleActiveProperty(tx, c, ts, req.Value)
		default:
			return &TradingSystemPropertyResponse{
				Status: ResponseStatusError,
				Message: "Unknown property : "+ req.Property,
			}, nil
	}
}

//=============================================================================

func handleRunningProperty(tx *gorm.DB, c *auth.Context, ts *db.TradingSystem, value string) (*TradingSystemPropertyResponse, error) {
	oldValue     := ts.Running
	newValue,err := getBool(value)
	if err != nil {
		return &TradingSystemPropertyResponse{
			Status : ResponseStatusError,
			Message: err.Error(),
		}, nil
	}

	if oldValue == newValue {
		return &TradingSystemPropertyResponse{
			Status : ResponseStatusSkipped,
		}, nil
	}


	if !oldValue && newValue {
		//--- Enabling
		c.Log.Info("handleRunningProperty: Starting trading system", "tsId", ts.Id)
	} else {
		//--- Disabling
		c.Log.Info("handleRunningProperty: Stopping trading system", "tsId", ts.Id)
	}

	ts.Running = newValue
	tradingsystem.UpdateStatus(ts)
	err = db.UpdateTradingSystem(tx, ts)
	if err != nil {
		return nil, err
	}

	err = updateRewind(ts)

	return &TradingSystemPropertyResponse{
		Status: ResponseStatusOk,
	}, err
}

//=============================================================================

func handleActivationProperty(tx *gorm.DB, c *auth.Context, ts *db.TradingSystem, value string) (*TradingSystemPropertyResponse, error) {
	oldValue     := ts.Activation
	newValue,err := getActivation(value)
	if err != nil {
		return &TradingSystemPropertyResponse{
			Status : ResponseStatusError,
			Message: err.Error(),
		}, nil
	}

	if oldValue == newValue {
		return &TradingSystemPropertyResponse{
			Status : ResponseStatusSkipped,
		}, nil
	}

	if oldValue == db.TsActivationAuto && newValue == db.TsActivationManual {
		//--- Activation = manual
		c.Log.Info("handleActivationProperty: Trading system's activation set to MANUAL", "tsId", ts.Id)
	} else {
		//--- Activation = auto
		c.Log.Info("handleActivationProperty: Trading system's activation set to AUTO", "tsId", ts.Id)
	}

	ts.Activation = newValue
	err = db.UpdateTradingSystem(tx, ts)

	return &TradingSystemPropertyResponse{
		Status: ResponseStatusOk,
	}, err
}

//=============================================================================

func handleActiveProperty(tx *gorm.DB, c *auth.Context, ts *db.TradingSystem, value string) (*TradingSystemPropertyResponse, error) {
	oldValue     := ts.Active
	newValue,err := getBool(value)
	if err != nil {
		return &TradingSystemPropertyResponse{
			Status : ResponseStatusError,
			Message: err.Error(),
		}, nil
	}

	if oldValue == newValue {
		return &TradingSystemPropertyResponse{
			Status : ResponseStatusSkipped,
		}, nil
	}

	if ts.Activation == db.TsActivationAuto {
		return &TradingSystemPropertyResponse{
			Status : ResponseStatusError,
			Message: "Trading system is in AUTOMATIC mode. Switch to MANUAL to pause",
		}, nil
	}

	if !oldValue && newValue {
		//--- Activating
		c.Log.Info("handleActiveProperty: Activating trading system", "tsId", ts.Id)
	} else {
		//--- Deactivating
		c.Log.Info("handleActiveProperty: Pausing trading system", "tsId", ts.Id)
	}

	ts.Active = newValue
	tradingsystem.UpdateStatus(ts)
	err = db.UpdateTradingSystem(tx, ts)
	if err != nil {
		return nil, err
	}

	err = updateRewind(ts)

	return &TradingSystemPropertyResponse{
		Status: ResponseStatusOk,
	}, err
}

//=============================================================================

func getBool(value string) (bool, error) {
	if value == "false" || value == "off" {
		return false, nil
	}
	if value == "true" || value == "on" {
		return true, nil
	}

	return false, errors.New("Unknown boolean value : "+ value)
}

//=============================================================================

func getActivation(value string) (db.TsActivation, error) {
	if value == "manual" {
		return db.TsActivationManual, nil
	}
	if value == "auto" {
		return db.TsActivationAuto, nil
	}

	return 0, errors.New("Unknown activation value : "+ value)
}

//=============================================================================

func updateRewind(ts *db.TradingSystem) error {
	return nil
}

//=============================================================================
