package main

import (
	"fmt"
	"net/http"
)

func handle(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<h1>Hello Golang</h1>")
}

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8000", nil)
}
