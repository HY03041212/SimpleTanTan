package main

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	//创建路由
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}