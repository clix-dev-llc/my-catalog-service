package main

import (
	"fmt"
	"github.com/rancher/go-rancher/client"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var Catalog map[string]Template

const catalogRoot string = "./DATA/templates/"

var metadataFolder = regexp.MustCompile(`^DATA/templates/[^/]+$`)

/*var Schemas = Schemas{
	Data: []Schema{
		{},
		{},
	},
}*/

type TemplateCollection struct {
	client.Collection
	Data []Template `json:"data,omitempty"`
}

func Init() {
	_, err := os.Stat(catalogRoot)
	if err != nil {
		cloneCatalog()
	}

	Catalog = make(map[string]Template)
	filepath.Walk(catalogRoot, walkCatalog)
}

func cloneCatalog() {
	fmt.Println("Cloning the catalog from github")
	//git clone the github repo
	e := exec.Command("git", "clone", "https://github.com/prachidamle/rancher-catalog", "./DATA")
	e.Run()
}

func walkCatalog(path string, f os.FileInfo, err error) error {

	if f.IsDir() && metadataFolder.MatchString(path) {
		fmt.Printf("Found a metadata folder, name:%s \n", f.Name())
		newTemplate := Template{}

		//read the root level config.yml

		filename, _ := filepath.Abs(path + "/config.yml")
		yamlFile, _ := ioutil.ReadFile(filename)
		config := make(map[string]string)
		_ = yaml.Unmarshal(yamlFile, &config)

		newTemplate.Name = config["name"]
		newTemplate.Category = config["category"]
		newTemplate.Description = config["description"]
		newTemplate.DefaultVersion = config["defaultVersion"]
		newTemplate.Path = f.Name() 

		//list the folders under the root level
		newTemplate.VersionLinks = make(map[string]string)
		dirList, _ := ioutil.ReadDir(path)
		for _, subfile := range dirList {
			if subfile.IsDir() {
				newTemplate.VersionLinks[subfile.Name()] = f.Name() + "/" + subfile.Name()
			}else if strings.HasPrefix(subfile.Name(), "catalogIcon"){
				newTemplate.IconLink = f.Name() + "/" + subfile.Name()
			}
		}
		Catalog[f.Name()] = newTemplate
	} 
	return nil
}

func readTemplateVersion(path string) Template{
	
	dirList, error := ioutil.ReadDir(catalogRoot + path)
	newTemplate := Template{}
	if error != nil{
		//return 404
	}else {
		for _, subfile := range dirList {
			if strings.HasPrefix(subfile.Name(), "config.yml") {
				filename, _ := filepath.Abs(catalogRoot + path + "/config.yml")
				yamlFile, _ := ioutil.ReadFile(filename)
				config := make(map[string]string)
				_ = yaml.Unmarshal(yamlFile, &config)
		
				newTemplate.Name = config["name"]
				newTemplate.Category = config["category"]
				newTemplate.Description = config["description"]
				newTemplate.DefaultVersion = config["defaultVersion"]
				newTemplate.Path = path 
				
			}else if strings.HasPrefix(subfile.Name(), "catalogIcon"){
				newTemplate.IconLink = path + "/" + subfile.Name()
			}else if strings.HasPrefix(subfile.Name(), "docker-compose"){
				filename, _ := filepath.Abs(catalogRoot + path + "/" + subfile.Name())
				composefile, _  := ioutil.ReadFile(filename)
				newTemplate.DockerCompose = string(composefile)
				
			}else if strings.HasPrefix(subfile.Name(), "rancher-compose"){
				filename, _ := filepath.Abs(catalogRoot + path + "/" + subfile.Name())
				composeBytes, _  := ioutil.ReadFile(filename)
				newTemplate.RancherCompose = string(composeBytes)
				
				//read the questions section
				
				RC := make(map[string]RancherCompose)
				_ = yaml.Unmarshal(composeBytes, &RC)
				
				newTemplate.Questions = RC["myService"].Questions
				
			}
		}
	}
	
	return newTemplate
	
}
