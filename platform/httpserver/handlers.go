package httpserver

import (
	"fmt"
	"net/http"
)

// HealthcheckHandler determines if the server is healthy.
func HealthcheckHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "OK")
}
