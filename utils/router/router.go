package router

import (
	"net/http"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(mux *http.ServeMux) *Router {
	return &Router{mux: mux}
}

// Http Methods
func (r *Router) Post(path string, handlerFunc http.HandlerFunc) {
	r.mux.Handle("POST "+path, handlerFunc)
}

func (r *Router) Get(path string, handlerFunc http.HandlerFunc) {
	r.mux.Handle("GET "+path, handlerFunc)
}

func (r *Router) Put(path string, handlerFunc http.HandlerFunc) {
	r.mux.HandleFunc("PUT "+path, handlerFunc)
}

func (r *Router) Delete(path string, handlerFunc http.HandlerFunc) {
	r.mux.HandleFunc("DELETE "+path, handlerFunc)
}

// ServeHTTP handles HTTP requests
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
