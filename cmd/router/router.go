package router

import (
	"github.com/charmbracelet/log"
	"net/http"
)

func SetUpRouter(proxys map[string]http.Handler) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler, ok := proxys[r.Host]
		log.Debug(r.Host, "proto", r.Proto)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		handler.ServeHTTP(w, r)
	})

	http.HandleFunc("OPTION /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
	})
}
