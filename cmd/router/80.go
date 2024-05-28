package router

import (
	"fmt"
	"github.com/charmbracelet/log"
	"net/http"
)

func ListenAndServe(proxymap map[string]http.Handler) {
	addr := fmt.Sprintf("0.0.0.0:%d", 80)
	log.Infof("Start listening %s", addr)
	mux := http.NewServeMux()
	SetUpCommonRouter(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler, ok := proxymap[r.Host]
		log.Debug(r.Host, "proto", r.Proto)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		handler.ServeHTTP(w, r)
	})

	log.Warn(http.ListenAndServe(addr, CorsMiddlewate(mux)))
}
