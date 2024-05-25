package middlewate

import "net/http"

type Middlewate func(handler http.Handler) http.Handler

func middlewateList() []Middlewate {
	return []Middlewate{corsMiddlewate}
}

func SetUpMiddlewate(proxys map[string]http.Handler) {
	middlewates := middlewateList()
	for s, handler := range proxys {
		for _, middlewate := range middlewates {
			proxys[s] = middlewate(handler)
		}
	}
}
