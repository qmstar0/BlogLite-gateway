package router

import (
	"crypto/tls"
	"fmt"
	"github.com/charmbracelet/log"
	"net/http"
)

func ListenAndServeTLS(proxymap map[string]http.Handler, certFp, keyFp string) {
	addr := fmt.Sprintf("0.0.0.0:%d", 443)
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

	server := http.Server{
		Addr:    ":443",
		Handler: mux,
		TLSConfig: &tls.Config{
			// 设置所需的加密套件
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			},
			MinVersion: tls.VersionTLS12,
		},
	}
	log.Warn(server.ListenAndServeTLS(certFp, keyFp))
}
