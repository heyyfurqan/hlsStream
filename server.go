package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	hlsDirectory := "."

	fileServer := http.FileServer(http.Dir(hlsDirectory))

	http.Handle("/", fileServer)

	port := 8080

	fmt.Printf("Serving %s on HTTP port: %d\n", hlsDirectory, port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
