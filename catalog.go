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
//	"github.com/rancher/rancher-compose/rancher"
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
				
				/*config := make(map[string]RancherCompose)
				//RC := RancherCompose{}
				_ = yaml.Unmarshal(composeBytes, &config)
				
				CS := config["myService"]
				
				fmt.Printf("rancherCompose.scale, name:%d \n", CS.Scale)
				fmt.Printf("rancherCompose.scale, name:%v \n", CS.Questions)
				//fmt.Printf("rancherCompose.scale, name:%s \n", config["myService"])
				/*for _, value := range RancherCompose.Questions {
					fmt.Printf("rancherCompose.scale, name:%s \n", value.Name)	
				}*/
				
				FirstQ := Question{}
				FirstQ.Name = "What is Question1?"
				FirstQ.Description = "This is Question1"
				FirstQ.Type = "string"
				
				SecondQ := Question{}
				SecondQ.Name = "What is Question2?"
				SecondQ.Description = "This is Question2"
				SecondQ.Type = "int"
				
				ThirdQ := Question{}
				ThirdQ.Name = "What is Question3?"
				ThirdQ.Description = "This is Question3"
				ThirdQ.Type = "boolean"
				
				FourthQ := Question{}
				FourthQ.Name = "What is Question4?"
				FourthQ.Description = "This is Question4"
				FourthQ.Type = "enum"
				FourthQ.Options = []string{"option1","option2"}
				
				newTemplate.Questions = append(newTemplate.Questions, FirstQ)
				newTemplate.Questions = append(newTemplate.Questions, SecondQ)
				newTemplate.Questions = append(newTemplate.Questions, ThirdQ)
				newTemplate.Questions = append(newTemplate.Questions, FourthQ)
				
			}
		}
	}
	
	return newTemplate
	
}
