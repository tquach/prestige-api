package metrics

import (
	"log"
	"net/http"
	"time"

	"github.com/tquach/prestige-api/platform/logger"
	"github.com/creack/dogstatsd"
)

// Stats holds the client connections to statsd server.
type Stats struct {
	statsdClient *dogstatsd.Client
	logger       logger.Logger
}

// Handler constructs a middleware handler to records metrics before and after a request is serviced.
func (s *Stats) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		duration := time.Since(start)
		s.logger.Debugf("\"%s %s\" - request took %s", req.Method, req.URL, duration)

		if s.statsdClient == nil {
			s.logger.Warning("Failed to report metrics; statsdClient is nil")
			return
		}

		s.statsdClient.Count("requests", 1, nil, 1)
		s.statsdClient.Histogram("requests.duration", duration.Seconds(), nil, 1)
	})
}

// Default creates a default Stats middleware component.
func Default(prefix string) *Stats {
	c, err := dogstatsd.New("127.0.0.1:8125")
	if err != nil {
		log.Println("WARN: Failed to initialize statsd client.")
	}

	return &Stats{
		statsdClient: c,
		logger:       logger.New(prefix),
	}
}
