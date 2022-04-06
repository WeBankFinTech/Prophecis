package restapi

import "net/http"

// GetHealth returns the health of the service
func GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
