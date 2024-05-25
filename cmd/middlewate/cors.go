package middlewate

import (
	"github.com/charmbracelet/log"
	"net/http"
)

func corsMiddlewate(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "http://blog.localhost")
		handler.ServeHTTP(w, r)
		log.Debug(w.Header(), "host", r.Host)
	})
}
