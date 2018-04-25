package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "grasshopper service: %s!", r.URL.Path[1:])

}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("Server is running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
