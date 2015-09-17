package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

func ListTemplates(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Request URL:, %q\n", html.EscapeString(r.URL.Path))

	//read the cloned repo starting at ./DATA/templates

    for _, value := range Catalog{
	  json.NewEncoder(w).Encode(value)
	}
}

func LoadTemplateDetails(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request URL:, %q", html.EscapeString(r.URL.Path))
}
