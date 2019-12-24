package accounts

import (
	"github.com/tquach/prestige-api/platform/app"
	"github.com/tquach/prestige-api/platform/auth"
	pg "gopkg.in/pg.v4"
)

// NewService creates a new web service to handle authentication requests.
func NewService(conn *pg.DB, tokenService auth.TokenService) app.Service {
	ctrl := NewController(conn, tokenService)

	authSvc := app.NewHTTPService("account-service")
	authSvc.AddRoute("POST", "/accounts/register/", ctrl.RegisterAccount)
	return authSvc
}
