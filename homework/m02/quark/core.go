package quark

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Engine struct {
	Routers
}

func New() *Engine {
	handlerMap := make(map[string]map[string]http.Handler)
	handlerMap["/"] = map[string]http.Handler{http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, "hello world")
	})}
	handlerMap["/metrics"] = map[string]http.Handler{http.MethodGet: promhttp.Handler()}

	return &Engine{
		Routers{
			m:          nil,
			handlerMap: handlerMap,
		},
	}
}

func notFoundErr(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprintln(w, "page not found")
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := e.handlerMap[r.URL.Path][r.Method]
	if !ok {
		handler = http.HandlerFunc(notFoundErr)
	}
	ms := Chain(e.m...)
	handler = ms(handler)
	handler.ServeHTTP(w, r)
}
