package main

import (
	mux "github.com/julienschmidt/httprouter"
)

type Route struct {
	Method string
	Path   string
	Handle mux.Handle
}

type Routes []Route

var routes = Routes{
	Route{
		"GET",
		"/",
		Index,
	},
	Route{
		"GET",
		"/posts",
		PostIndex,
	},
	Route{
		"GET",
		"/posts/:id",
		PostShow,
	},
	Route{
		"POST",
		"/posts",
		PostCreate,
	},
}

func NewRouter() *mux.Router {
	router := mux.New()
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.Handle)
	}
	return router
}
