package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"mixloop/handlers"
	"mixloop/utils"
)

func main() {
	// Create uploads and output directories
	os.MkdirAll("uploads", 0755)
	os.MkdirAll("output", 0755)

	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/api/mix", handlers.MixAudioHandler).Methods("POST")
	r.HandleFunc("/api/progress", utils.ProgressHandler).Methods("GET")
	r.HandleFunc("/ws/progress", utils.WebSocketHandler)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// CORS - Allow all origins for web/VPS deployment
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	handler := c.Handler(r)

	log.Println("MixLoop Audio API Server - Developed by BITZY.ID")
	log.Println("Server starting on port 8081")
	log.Fatal(http.ListenAndServe(":8081", handler))
}
