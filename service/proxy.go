package service

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(targetPath string) (http.Handler, error) {
	target, err := url.Parse(targetPath)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		if req.TLS != nil {
			req.Header.Set("X-Forwarded-Proto", "HTTPS")
		} else {
			req.Header.Set("X-Forwarded-Proto", "HTTP")
		}
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-Forwarded-From", req.RemoteAddr)
	}
	return proxy, nil
}
