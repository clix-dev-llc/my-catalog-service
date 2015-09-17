package main

import (
	"fmt"
	"os/exec"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var Catalog map[string]Template

func Init(){
	_, err := os.Stat("./DATA/templates")
    if err != nil {
    // check if the source is indeed a directory or not
		cloneCatalog()
    }
    
    Catalog = make(map[string]Template)
    
    var root = "./DATA/templates"

	filepath.Walk(root, walkCatalog)
}

func cloneCatalog(){
	 fmt.Println("Cloning the catalog from github")
	//git clone the github repo
	e := exec.Command("git", "clone", "https://github.com/prachidamle/rancher-catalog", "./DATA");
	e.Run()
}

func walkCatalog(path string, f os.FileInfo, err error) error {
	fmt.Printf("Walking path %s with file name %s \n", path, f.Name())

	if f.IsDir() {
		fmt.Printf("Found a template name:%s \n", f.Name())
		newTemplate := Template{}
		newTemplate.Name = f.Name()
		Catalog[f.Name()] = newTemplate
	} else if strings.HasSuffix(f.Name(), "config.yml") {
		filename, _ := filepath.Abs(path)
		yamlFile, _ := ioutil.ReadFile(filename)
		config := make(map[string]string)
		_ = yaml.Unmarshal(yamlFile, &config)
		t := Catalog[config["name"]]
		t.Category = config["category"]
		t.Description = config["description"]
		t.DefaultVersion = config["defaultVersion"]
		Catalog[config["name"]] = t
	}

	return nil
}