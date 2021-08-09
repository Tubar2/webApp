package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router is a mux router with extra features
type Router struct {
	*mux.Router
}

func (r *Router) GET(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	return r.HandleFunc(path, f).Methods("GET")
}

func (r *Router) POST(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	return r.HandleFunc(path, f).Methods("POST")

}

func (r *Router) PUT(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	return r.HandleFunc(path, f).Methods("PUT")

}

func (r *Router) DELETE(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	return r.HandleFunc(path, f).Methods("DELETE")

}

func NewRouter() *Router {
	return &Router{
		Router: mux.NewRouter(),
	}
}
