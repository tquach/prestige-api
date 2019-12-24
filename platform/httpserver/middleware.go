package httpserver

import "net/http"

// Middleware defines a function that is applied before or after a request is serviced.
type Middleware func(http.Handler) http.Handler
