package main

import (
	"log"
	"net/http"
)

func main() {
	Init()

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8088", router))
}
