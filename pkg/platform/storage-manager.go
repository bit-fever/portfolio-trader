//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

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

package platform

import (
	"bytes"
	"context"
	"github.com/bit-fever/core"
	"github.com/bit-fever/core/req"
	"github.com/coreos/go-oidc/v3/oidc"
	"log/slog"
	"net/http"
	"strconv"
)

//=============================================================================
type EquityRequest struct {
	Image []byte `json:"image"`
}


func SetEquityChart(id uint, data []byte) error {
	slog.Info("SetEquityChart: Sending equity chart to storage manager", "id", id)

	client :=req.GetClient("bf")
//	url := c.Config.(*app.Config).Platform.Storage +"/v1/trading-systems/"+id+"/equity-chart"
	url := "https://bitfever-server:8452/api/storage/v1/trading-systems/"+strconv.Itoa(int(id))+"/equity-chart"
	par := EquityRequest{ Image: data}
	token := aoh()
	res := ""
	err := req.DoPost(client, url, &par, &res, token)

	if err != nil {
		slog.Error("SetEquityChart: Got an error when sending to storage-manager", "id", id, "error", err.Error())
		return req.NewServerError("Cannot communicate with storage-manager: %v", err.Error())
	}

	slog.Info("SetEquityChart: Equity chart saved", "id", id)
	return nil
}

//=============================================================================

func aoh() string {
	authority     := "https://bitfever-server:8443/auth/realms/bitfever"
	clientID      := "bitfever-backend"
	clientSecret  := "AXRNrPrBiiOlBcwiOPX67XFTcbu9p9IZ"

	client        := req.GetClient("bf")
	ccontext      := oidc.ClientContext(context.Background(), client)
	provider, err := oidc.NewProvider(ccontext, authority)
	core.ExitIfError(err)

	params := "grant_type=client_credentials&client_id="+clientID+"&client_secret="+ clientSecret
	result := TokenResponse{}
	url    := "https://bitfever-server:8443/auth/realms/bitfever/protocol/openid-connect/token"

	_ = myPost(client, url, params, &result)
	t := provider.Endpoint()
	slog.Info("Provider : "+provider.UserInfoEndpoint()+"| "+ t.AuthURL)


//	config := oauth2.Config{
//		ClientID:     clientID,
//		ClientSecret: clientSecret,
//		Endpoint:     provider.Endpoint(),
////		RedirectURL:  "http://127.0.0.1:5556/auth/google/callback",
//		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
//	}
//
//	token, err := config.PasswordCredentialsToken(ccontext, "", "")
	if err != nil {
		slog.Error("Error : "+ err.Error())
	} else {
//		slog.Info("Ok : "+ token.AccessToken)
	}
	return result.AccessToken
}

//=============================================================================

func myPost(client *http.Client, url string, params string, output any) error {
	body := []byte(params)
	reader := bytes.NewReader(body)

	rq, err := http.NewRequest("POST", url, reader)
	if err != nil {
		slog.Error("Error creating a POST request", "error", err.Error())
		return err
	}

	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(rq)
	return req.BuildResponse(res, err, &output)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}
