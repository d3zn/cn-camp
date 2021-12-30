package quark

import "net/http"

var _ Router = &Routers{}

type Router interface {
	GET(string, http.HandlerFunc) Router
	POST(string, http.HandlerFunc) Router
}

type Routers struct {
	m          []Middleware
	handlerMap map[string]map[string]http.Handler
}

func (r *Routers) Use(m ...Middleware) {
	r.m = append(r.m, m...)
}

func (r *Routers) GET(path string, handler http.HandlerFunc) Router {
	r.add(http.MethodGet, path, handler)
	return r
}

func (r *Routers) POST(path string, handler http.HandlerFunc) Router {
	r.add(http.MethodPost, path, handler)
	return r
}

func (r *Routers) add(method, path string, handler http.Handler) {
	r.handlerMap[path] = map[string]http.Handler{method: handler}
}
