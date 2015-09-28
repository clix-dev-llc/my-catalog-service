package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rancher/go-rancher/client"
	"net/http"
)

func ListTemplates(w http.ResponseWriter, r *http.Request) {
	//read the catalog

	resp := TemplateCollection{}
	for _, value := range Catalog {
		value.VersionLinks = PopulateTemplateLinks(r, &value)
		PopulateResource(r, "template", value.Name, &value.Resource)
		resp.Data = append(resp.Data, value)
	}

	PopulateCollection(&resp.Collection)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)

}

func LoadTemplateMetadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["templateId"]
	templateMetadata, ok := Catalog[path]
	if ok {
		templateMetadata.VersionLinks = PopulateTemplateLinks(r, &templateMetadata)
		PopulateResource(r, "template", templateMetadata.Name, &templateMetadata.Resource)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(templateMetadata)
	} else {
		//404
	}
}

func LoadTemplateVersion(w http.ResponseWriter, r *http.Request){
	//read the template version from disk
	vars := mux.Vars(r)	
	path := vars["templateId"] + "/" + vars["versionId"]
	template := readTemplateVersion(path)
	template.VersionLinks = PopulateTemplateLinks(r, &template)
	PopulateResource(r, "template", template.Name, &template.Resource)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(template)
}


func LoadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := "DATA/templates/" + vars["templateId"] + "/" + vars["versionId"] + "/" + vars["imageId"]
	fmt.Printf("path is%s.\n", path)
	http.ServeFile(w, r, path)
}

func PopulateCollection(collection *client.Collection) {
	collection.Type = "collection"
	collection.ResourceType = "template"
}

func PopulateTemplateLinks(r *http.Request, template *Template) map[string]string {

	copyOfversionLinks := make(map[string]string)
	for key, value := range template.VersionLinks {
		copyOfversionLinks[key] = BuildURL(r, "template", value)
	}

	template.IconLink = BuildURL(r, "image", template.IconLink)
	
	return copyOfversionLinks
}

func PopulateResource(r *http.Request, resourceType, resourceId string, resource *client.Resource) {
	resource.Type = resourceType

	selfLink := BuildURL(r, "template", resourceId)

	resource.Links = map[string]string{
		"self": selfLink,
	}
}

func BuildURL(r *http.Request, resourceType, resourceId string) string {

	var scheme string = "http://"
	var host string = r.Host
	var pluralName string = resourceType + "s"
	var version string = "v1"

	//get the baseURI
	return scheme + host + "/" + version + "/" + pluralName + "/" + resourceId

}
