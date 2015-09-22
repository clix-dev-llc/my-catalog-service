package main

import (
	"encoding/json"
	"github.com/rancher/go-rancher/client"
	"net/http"
	"github.com/gorilla/mux"
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

func LoadTemplateDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["templatePath"]
	templateMetadata, ok := Catalog[path]
	if ok{
		templateMetadata.VersionLinks = PopulateTemplateLinks(r, &templateMetadata)
		PopulateResource(r, "template", templateMetadata.Name, &templateMetadata.Resource)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(templateMetadata)
	}

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
    var pluralName string =  resourceType + "s"
    var version string = "v1"
    
    //get the baseURI
    return scheme + host + "/" + version + "/" + pluralName + "/" + resourceId
    

}
