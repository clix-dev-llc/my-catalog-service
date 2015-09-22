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
	//"strings"
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
		// check if the source is indeed a directory or not
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
		newTemplate.Name = f.Name()

		//read the root level config.yml

		filename, _ := filepath.Abs(path + "/config.yml")
		yamlFile, _ := ioutil.ReadFile(filename)
		config := make(map[string]string)
		_ = yaml.Unmarshal(yamlFile, &config)

		newTemplate.Category = config["category"]
		newTemplate.Description = config["description"]
		newTemplate.DefaultVersion = config["defaultVersion"]

		//list the folders under the root level
		newTemplate.VersionLinks = make(map[string]string)
		dirList, _ := ioutil.ReadDir(path)
		for _, subfile := range dirList {
			if subfile.IsDir() {
				newTemplate.VersionLinks[subfile.Name()] = f.Name() + "/" + subfile.Name()
			}
		}
		newTemplate.IconLink = "link to icon"

		Catalog[f.Name()] = newTemplate
	} else {
		fmt.Printf("Walking path %s with file name %s \n", path, f.Name())

	}

	return nil
}
