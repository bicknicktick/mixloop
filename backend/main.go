package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"mixloop/handlers"
)

func main() {
	// Create uploads and output directories
	os.MkdirAll("uploads", 0755)
	os.MkdirAll("output", 0755)

	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/mix", handlers.MixAudioHandler).Methods("POST")
	r.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	log.Println("MixLoop Audio API Server - Developed by BITZY.ID")
	log.Println("Server starting on port 8081")
	log.Fatal(http.ListenAndServe(":8081", handler))
}
