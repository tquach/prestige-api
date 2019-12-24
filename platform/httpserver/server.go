package httpserver

import (
	"fmt"
	"net/http"

	"github.com/tquach/prestige-api/platform/app"
	"github.com/tquach/prestige-api/platform/logger"
	_ "github.com/lib/pq" // postgresql
)

// Opts allows users to pass in configuration settings.
type Opts struct {
	Hostname string
	Port     int
}

// AppServer defines the application server construct.
type AppServer struct {
	options *Opts
	logger  logger.Logger
}

// Start will start the application server.
func (a *AppServer) Start() error {
	addr := fmt.Sprintf("%s:%d", a.options.Hostname, a.options.Port)
	a.logger.Infof("Starting server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

// RegisterService registers the service and it's routes.
func (a *AppServer) RegisterService(prefix string, svc app.Service) error {
	a.logger.Infof("Registering %q with service %q", prefix, svc.Name())
	http.Handle(prefix, svc.Router())
	return nil
}

// Shutdown attemptes to gracefully shut down the server and clean up resources.
func (a *AppServer) Shutdown() error {
	return nil
}

// NewAppServer will construct a new application server and initialize any default routes.
func NewAppServer(opts *Opts) (*AppServer, error) {
	appServer := &AppServer{
		options: opts,
		logger:  logger.New("prestige-api"),
	}

	// Register endpoints
	http.HandleFunc("/healthcheck", HealthcheckHandler)
	return appServer, nil
}
