package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		//路由名称、请求方法、url路径、处理方法
		"GetUsers",
		"GET",
		"/users",
		GetUsers,
	},
	Route{
		"UserCreate",
		"POST",
		"/users",
		UserCreate,
	},
	Route{
		"GetRelationships",
		"GET",
		"/users/{UserId}/relationships",
		GetRelationships,
	},
	Route{
		"PutRelationships",
		"PUT",
		"/users/{UserId}/relationships/{OtherUserId}",
		PutRelationships,
	},
}
