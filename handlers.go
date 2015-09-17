package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/rancher/go-rancher/client"
)

func ListTemplates(w http.ResponseWriter, r *http.Request) {

	//fmt.Fprintf(w, "Request URL:, %q\n", html.EscapeString(r.URL.Path))

	//read the cloned repo starting at ./DATA/templates

	resp := TemplateCollection{}
	for _, value := range Catalog {
		PopulateResource(&value.Resource)
		resp.Data = append(resp.Data, value)
	}

	PopulateCollection(&resp.Collection)
	json.NewEncoder(w).Encode(resp)
}

func LoadTemplateDetails(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request URL:, %q", html.EscapeString(r.URL.Path))
}

func PopulateCollection(collection *client.Collection) {
	collection.Type = "collection"
}

func PopulateResource(r *http.Request, resourceType, resourceId string, resource *client.Resource) {
	resource.Type = "template"
	resource.Links = map[string]string{
		"self": "foo",
	}
}
