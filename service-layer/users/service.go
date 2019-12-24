package users

import (
	"github.com/tquach/prestige-api/platform/app"
	"github.com/tquach/prestige-api/platform/auth"
	pg "gopkg.in/pg.v4"

	_ "github.com/jackc/pgx" // postgres driver
)

// NewService creates a new user service.
func NewService(db *pg.DB, tokenSvc auth.TokenService) app.Service {
	ctrl := NewController(db)
	userSvc := app.NewSecuredHTTPService("user-service", tokenSvc)

	userSvc.AddRoute("GET", "/users/:id", ctrl.Find)
	userSvc.AddRoute("GET", "/users/", ctrl.FindUserBySocialID)
	return userSvc
}
