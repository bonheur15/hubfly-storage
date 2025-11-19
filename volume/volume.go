package volume

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	log.Printf("Command: %s %v\nOutput: %s", name, args, output)
	if err != nil {
		return fmt.Errorf("%v: %s", err, output)
	}
	return nil
}

func CreateVolume(name, size, baseDir string) (string, error) {
	volumePath := filepath.Join(baseDir, name)
	dataPath := filepath.Join(volumePath, "_data")
	absDataPath, err := filepath.Abs(dataPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %v", err)
	}
	imagePath := filepath.Join(volumePath, "volume.img")

	if size == "" {
		size = "1G"
	}

	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	log.Printf("Allocating %s image file at %s", size, imagePath)
	if err := runCommand("sudo", "fallocate", "-l", size, imagePath); err != nil {
		return "", fmt.Errorf("fallocate failed: %v", err)
	}

	log.Printf("Formatting %s as ext4", imagePath)
	if err := runCommand("sudo", "mkfs.ext4", imagePath); err != nil {
		return "", fmt.Errorf("mkfs.ext4 failed: %v", err)
	}

	log.Printf("Mounting volume image at %s", dataPath)
	if err := runCommand("sudo", "mount", "-o", "loop", imagePath, dataPath); err != nil {
		return "", fmt.Errorf("mount failed: %v", err)
	}

	log.Printf("Registering docker volume: %s", name)
	if err := runCommand(
		"docker", "volume", "create",
		"--name", name,
		"--opt", fmt.Sprintf("device=%s", absDataPath),
		"--opt", "type=none",
		"--opt", "o=bind",
	); err != nil {
		return "", fmt.Errorf("docker volume create failed: %v", err)
	}

	return name, nil
}

func DeleteVolume(name, baseDir string) error {
	volumePath := filepath.Join(baseDir, name)
	dataPath := filepath.Join(volumePath, "_data")

	log.Printf("Unmounting volume at %s", dataPath)
	if err := runCommand("sudo", "umount", dataPath); err != nil {
		// Log the error but continue, as the volume might not be mounted
		log.Printf("unmount failed (might be acceptable if not mounted): %v", err)
	}

	log.Printf("Removing docker volume: %s", name)
	if err := runCommand("docker", "volume", "rm", name); err != nil {
		return fmt.Errorf("docker volume rm failed: %v", err)
	}

	log.Printf("Removing volume directory: %s", volumePath)
	if err := os.RemoveAll(volumePath); err != nil {
		return fmt.Errorf("failed to remove volume directory: %v", err)
	}

	return nil
}