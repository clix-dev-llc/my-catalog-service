package main

type Template struct {
	Name           string `json:"name"`
	Category       string `json:"category"`
	Description    string `json:"description"`
	DefaultVersion string `json:"version"`
}

type Templates []Template
