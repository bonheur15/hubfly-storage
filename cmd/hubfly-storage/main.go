package main

import (
	"log"
	"net/http"
	"os"

	"hubfly-storage/handlers"
)

func main() {
	baseDir := "./docker/volumes"
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Fatalf("Failed to create base directory: %v", err)
	}

	http.HandleFunc("/create-volume", handlers.CreateVolumeHandler(baseDir))
	http.HandleFunc("/delete-volume", handlers.DeleteVolumeHandler(baseDir))
	http.HandleFunc("/health", handlers.HealthCheckHandler())
	        http.HandleFunc("/volume-stats", handlers.GetVolumeStatsHandler(baseDir))
	        http.HandleFunc("/volumes", handlers.GetVolumesHandler(baseDir))
        http.HandleFunc("/volumes/stats", handlers.GetVolumesStatsHandler(baseDir))
	log.Println("🚀 Server running on port 8203...")
	if err := http.ListenAndServe(":8203", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
