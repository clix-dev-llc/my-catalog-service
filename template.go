package main

import "github.com/rancher/go-rancher/client"

type Template struct {
	client.Resource
	Name           string            `json:"name"`
	Category       string            `json:"category"`
	Description    string            `json:"description"`
	DefaultVersion string            `json:"defaultVersion"`
	IconLink       string            `json:"iconLink"`
	VersionLinks   map[string]string `json:"versionLinks"`
	DockerCompose  string            `json:"dockerCompose"`
	RancherCompose string            `json:"rancherCompose"`
	Questions      []Question	     `json:"questions"`
}

var Templates map[string]Template

type Question struct {
	Name           string            `yaml:"name"`
	Description    string            `yaml:"description"`
	Type		   string            `yaml:"type"`
	Options		   []string          `yaml:"options"`
}