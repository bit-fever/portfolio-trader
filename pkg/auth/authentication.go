package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//=============================================================================

type UnauthorizedResponse struct {
	Status   string `json:"status"   example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"401"`
	Message  string `json:"message"  example:"Authorisation failed"`
	Details  string `json:"details"  example:"..."`
}

//=============================================================================

//claims component of jwt contains mainy fields , we need only roles of DemoServiceClient
//"DemoServiceClient":{"DemoServiceClient":{"roles":["pets-admin","pet-details","pets-search"]}},

type Claims struct {
	ResourceAccess client `json:"resource_access,omitempty"`
	JTI            string `json:"jti,omitempty"`
}

type client struct {
	DemoServiceClient clientRoles `json:"DemoServiceClient,omitempty"`
}

type clientRoles struct {
	Roles []string `json:"roles,omitempty"`
}

var RealmConfigURL string = "http://localhost:8200/auth/realms/bitfever"
var clientID string = "bf-frontend"

//=============================================================================

func Wrap(handler gin.HandlerFunc, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawAccessToken := c.GetHeader("Authorization")

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{
			Timeout:   time.Duration(6000) * time.Second,
			Transport: tr,
		}
		ctx := oidc.ClientContext(context.Background(), client)
		provider, err := oidc.NewProvider(ctx, RealmConfigURL)
		if err != nil {
			authorisationFailed(c, "Cannot get the provider", err)
			return
		}

		oidcConfig := &oidc.Config{
			ClientID: clientID,
		}
		verifier := provider.Verifier(oidcConfig)
		idToken, err := verifier.Verify(ctx, rawAccessToken)
		if err != nil {
			authorisationFailed(c, "Failed to verify the token", err)
			return
		}

		var idTokenClaims Claims
		if err := idToken.Claims(&idTokenClaims); err != nil {
			authorisationFailed(c, "Failed to get claims", err)
			return
		}

		fmt.Println(idTokenClaims)
		//checking the roles
		user_access_roles := idTokenClaims.ResourceAccess.DemoServiceClient.Roles
		for _, b := range user_access_roles {
			if b == role {
				handler(c)
				return
			}
		}

		authorisationFailed(c, "User not allowed to access this api", nil)
	}
}

//=============================================================================

func authorisationFailed(c *gin.Context, message string, err error) {

	details := ""
	if err != nil {
		details = err.Error()
	}

	c.JSON(http.StatusUnauthorized, UnauthorizedResponse{
		Status:   "FAILED",
		HTTPCode: http.StatusUnauthorized,
		Message:  message,
		Details:  details,
	})
}

//=============================================================================
