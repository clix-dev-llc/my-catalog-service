package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

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

var routes = Routes{
	Route{
		"ListTemplates",
		"GET",
		"/v1/templates",
		ListTemplates,
	},
	Route{
		"LoadTemplateDetails",
		"GET",
		"/v1/templates/{templateId}",
		LoadTemplateMetadata,
	},
	Route{
		"LoadTemplateDetails",
		"GET",
		"/v1/templates/{templateId}/{versionId}",
		LoadTemplateVersion,
	},
	Route{
		"LoadVersionImage",
		"GET",
		"/v1/images/{templateId}/{versionId}/{imageId}",
		LoadImage,
	},
	Route{
		"LoadImage",
		"GET",
		"/v1/images/{templateId}/{imageId}",
		LoadImage,
	},
}
