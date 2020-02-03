package main

import (
	"log"
	"net/http"
	functions "functions"
	"fmt"
)

func main() {
	fmt.Println("Setting up test server.")

	http.HandleFunc("/", functions.HelloWorld)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
