package app

import (
	"net/http"

	"github.com/tquach/prestige-api/platform/auth"
	"github.com/tquach/prestige-api/platform/config"
	"github.com/tquach/prestige-api/platform/logger"
	"github.com/tquach/prestige-api/platform/metrics"
	"github.com/bmizerany/pat"
	"github.com/rs/cors"
)

// Config defines service configuration settings.
type Config struct {
	DatabaseURL string
}

// Service describes a service for implementation. A Service implicitly extends http.Handler.
type Service interface {
	// Name returns the name of this service.
	Name() string
	Router() http.Handler
}

// Middleware defines a function that is applied before or after a request is serviced.
type Middleware func(http.Handler) http.Handler

// HTTPService has common HTTP functionality for standard web services.
type HTTPService struct {
	name       string
	secured    bool
	debug      bool
	mux        *pat.PatternServeMux
	logger     logger.Logger
	middleware []Middleware
}

// SecureHTTPService is a HTTPService variant that requires valid a access token for each request.
type SecureHTTPService struct {
	*HTTPService
	tokenSvc auth.TokenService
}

// Use adds an additional middleware function to the chain.
func (h *HTTPService) Use(m Middleware) {
	h.middleware = append(h.middleware, m)
}

// Name returns the name of the service
func (h *HTTPService) Name() string {
	return h.name
}

// AddRoute will map a handler to an endpoint and HTTP method. All HTTP methods supported.
func (h *HTTPService) AddRoute(method string, path string, hdlr http.HandlerFunc) {
	h.mux.Add(method, path, http.HandlerFunc(hdlr))
}

// Router returns the handler with all middleware associated.
func (h *HTTPService) Router() http.Handler {
	corsHandler := cors.New(cors.Options{
		Debug:          h.debug,
		AllowedMethods: []string{"GET", "POST", "PUT"},
		AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "Authorization"},
	})
	metricsHandler := metrics.Default(h.Name())
	handler := http.Handler(h.mux)
	for i := len(h.middleware) - 1; i >= 0; i-- {
		handler = h.middleware[i](handler)
	}
	return corsHandler.Handler(metricsHandler.Handler(handler))
}

// Router returns a secured router.
func (s *SecureHTTPService) Router() http.Handler {
	r := s.HTTPService.Router()
	securedHandler := auth.SecuredHandler(s.tokenSvc)
	return securedHandler.Handler(r)
}

// NewHTTPService creates an instance of HTTPService and initializes internal fields.
func NewHTTPService(name string) *HTTPService {
	debug := config.Env == config.Local || config.Env == config.Dev
	return &HTTPService{
		mux:    pat.New(),
		name:   name,
		logger: logger.New(name),
		debug:  debug,
	}
}

// NewSecuredHTTPService creates an instance of HTTPService and initializes internal fields.
func NewSecuredHTTPService(name string, tokenSvc auth.TokenService) *SecureHTTPService {
	return &SecureHTTPService{
		HTTPService: NewHTTPService(name),
		tokenSvc:    tokenSvc,
	}
}
