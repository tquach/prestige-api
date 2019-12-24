package trips

import (
	"github.com/tquach/prestige-api/platform/app"
	"github.com/tquach/prestige-api/platform/auth"
	_ "github.com/jackc/pgx" // postgres driver
	pg "gopkg.in/pg.v4"
)

// NewService creates a new service for managing trips.
func NewService(conn *pg.DB, tokenSvc auth.TokenService) app.Service {
	ctrl := NewController(conn)
	tripSvc := app.NewSecuredHTTPService("trip-service", tokenSvc)

	tripSvc.AddRoute("POST", "/trips/ideas/", ctrl.SaveIdea)
	// tripSvc.AddRoute("PUT", "/trips/ideas/", ctrl.UpdateIdea)

	tripSvc.AddRoute("POST", "/trips/destinations/", ctrl.SaveDestination)
	tripSvc.AddRoute("PUT", "/trips/destinations/:id", ctrl.UpdateDestination)

	tripSvc.AddRoute("GET", "/trips/:id", ctrl.Find)
	// tripSvc.AddRoute("PUT", "/trips/:id", ctrl.Update)

	tripSvc.AddRoute("POST", "/trips/", ctrl.Save)
	tripSvc.AddRoute("GET", "/trips/", ctrl.ListTrips)

	// tripSvc.AddRoute("POST", "/ideas/", ctrl.)
	return tripSvc
}
