package main

import (
	"fmt"
	"log"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func main() {
	http.HandleFunc("/ping", pingHandler)

	port := ":8080"
	fmt.Printf("Service statring on port %s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Failed to start server")
	}
}
