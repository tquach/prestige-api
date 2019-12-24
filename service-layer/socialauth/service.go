package socialauth

import (
	"github.com/tquach/prestige-api/platform/app"
	"github.com/tquach/prestige-api/platform/auth"
	pg "gopkg.in/pg.v4"
)

// NewService creates a new web service to handle authentication requests.
func NewService(conn *pg.DB, tokenSvc auth.TokenService) app.Service {
	ctrl := NewController(tokenSvc, conn)

	authSvc := app.NewHTTPService("auth-service")
	authSvc.AddRoute("POST", "/auth/", ctrl.Authenticate)
	return authSvc
}
