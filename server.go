package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func StartServer() {
	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %s", err)
	}

	fileServer := http.FileServer(http.Dir(currentDirectory))

	http.Handle("/", corsMiddleware(fileServer))

	port := 8080

	fmt.Printf("Serving on HTTP port: %d\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
