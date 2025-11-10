package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type DockerVolumePayload struct {
	Name       string            `json:"Name"`
	Driver     string            `json:"Driver"`
	DriverOpts map[string]string `json:"DriverOpts"`
	Labels     map[string]string `json:"Labels"`
}

func main() {
	baseDir := "./docker/volumes"
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Fatalf("Failed to create base directory: %v", err)
	}

	http.HandleFunc("/create-volume", func(w http.ResponseWriter, r *http.Request) {
		var payload DockerVolumePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON payload: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("Received request to create volume: %s", payload.Name)

		volumePath := filepath.Join(baseDir, payload.Name)
		dataPath := filepath.Join(volumePath, "_data")
		imagePath := filepath.Join(volumePath, "volume.img")

		size := payload.DriverOpts["size"]
		if size == "" {
			size = "1G"
		}

		if err := os.MkdirAll(dataPath, 0755); err != nil {
			handleError(w, fmt.Sprintf("Failed to create directory: %v", err))
			return
		}

		log.Printf("Allocating %s image file at %s", size, imagePath)
		if err := runCommand("sudo", "fallocate", "-l", size, imagePath); err != nil {
			handleError(w, fmt.Sprintf("fallocate failed: %v", err))
			return
		}

		log.Printf("Formatting %s as ext4", imagePath)
		if err := runCommand("sudo", "mkfs.ext4", imagePath); err != nil {
			handleError(w, fmt.Sprintf("mkfs.ext4 failed: %v", err))
			return
		}

		log.Printf("Mounting volume image at %s", dataPath)
		if err := runCommand("sudo", "mount", "-o", "loop", imagePath, dataPath); err != nil {
			handleError(w, fmt.Sprintf("mount failed: %v", err))
			return
		}

		log.Printf("Registering docker volume: %s", payload.Name)
		if err := runCommand(
			"docker", "volume", "create",
			"--name", payload.Name,
			"--opt", fmt.Sprintf("device=%s", dataPath),
			"--opt", "type=none",
			"--opt", "o=bind",
		); err != nil {
			handleError(w, fmt.Sprintf("docker volume create failed: %v", err))
			return
		}

		log.Printf("Volume %s created successfully!", payload.Name)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
			"name":   payload.Name,
		})
	})

	log.Println("🚀 Server running on port 8203...")
	if err := http.ListenAndServe(":8203", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	log.Printf("Command: %s %v\nOutput: %s", name, args, output)
	if err != nil {
		return fmt.Errorf("%v: %s", err, output)
	}
	return nil
}

func handleError(w http.ResponseWriter, msg string) {
	log.Println("❌ " + msg)
	http.Error(w, msg, http.StatusInternalServerError)
}
