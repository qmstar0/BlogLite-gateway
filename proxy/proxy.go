package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(sourcePath, targetPath string) (string, http.Handler, error) {
	source, err := url.Parse(sourcePath)
	if err != nil {
		return "", nil, err
	}

	target, err := url.Parse(targetPath)
	if err != nil {
		return "", nil, err
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
	return source.Host, proxy, nil
}
