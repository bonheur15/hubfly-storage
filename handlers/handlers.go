package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"hubfly-storage/volume"
)

type DockerVolumePayload struct {
	Name       string            `json:"Name"`
	Driver     string            `json:"Driver"`
	DriverOpts map[string]string `json:"DriverOpts"`
	Labels     map[string]string `json:"Labels"`
}

func handleError(w http.ResponseWriter, msg string, statusCode int) {
	log.Println("❌ " + msg)
	http.Error(w, msg, statusCode)
}

func CreateVolumeHandler(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload DockerVolumePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			handleError(w, fmt.Sprintf("Invalid JSON payload: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("Received request to create volume: %s", payload.Name)

		size := payload.DriverOpts["size"]
		if size == "" {
			size = "1G"
		}

		volName, err := volume.CreateVolume(payload.Name, size, baseDir)
		if err != nil {
			handleError(w, fmt.Sprintf("Failed to create volume: %v", err), http.StatusInternalServerError)
			return
		}

		log.Printf("Volume %s created successfully!", volName)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
			"name":   volName,
		})
	}
}

func DeleteVolumeHandler(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload DockerVolumePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			handleError(w, fmt.Sprintf("Invalid JSON payload: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("Received request to delete volume: %s", payload.Name)

		if err := volume.DeleteVolume(payload.Name, baseDir); err != nil {
			handleError(w, fmt.Sprintf("Failed to delete volume: %v", err), http.StatusInternalServerError)
			return
		}

		log.Printf("Volume %s deleted successfully!", payload.Name)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
			"name":   payload.Name,
		})
	}
}

func HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}
}