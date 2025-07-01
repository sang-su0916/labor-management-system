package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// 가장 단순한 HTTP 서버
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Labor Management System!")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"healthy","service":"labor-management-system"}`)
	})

	addr := "0.0.0.0:" + port
	log.Printf("Simple server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}