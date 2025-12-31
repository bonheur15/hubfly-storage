package main

import (
	"log"
	"net/http"

	"hubfly-storage/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	baseDir := "./docker-volumes"

	http.HandleFunc("/create-volume", handlers.CreateVolumeHandler(baseDir))
	http.HandleFunc("/delete-volume", handlers.DeleteVolumeHandler(baseDir))
	http.HandleFunc("/health", handlers.HealthCheckHandler())
	http.HandleFunc("/volume-stats", handlers.GetVolumeStatsHandler(baseDir))
	http.HandleFunc("/dev/volumes", handlers.GetVolumesHandler(baseDir))
	http.HandleFunc("/url-volume/create", handlers.URLVolumeCreateHandler)

	log.Println("🚀 Server running on port 8203...")
	if err := http.ListenAndServe(":8203", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
