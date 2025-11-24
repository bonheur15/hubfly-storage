package volume

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type VolumeStats struct {
	Name      string `json:"name"`
	Size      string `json:"size"`
	Used      string `json:"used"`
	Available string `json:"available"`
	Usage     string `json:"usage"`
	MountPath string `json:"mount_path"`
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

func runCommandWithOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	log.Printf("Command: %s %v\nOutput: %s", name, args, output)
	if err != nil {
		return string(output), fmt.Errorf("%v: %s", err, output)
	}
	return string(output), nil
}

// volumeExists checks if a docker volume with the given name already exists.
func volumeExists(name string) (bool, error) {
	output, err := runCommandWithOutput("docker", "volume", "ls", "-q", "-f", "name="+name)
	if err != nil {
		return false, fmt.Errorf("failed to check if volume exists: %v", err)
	}
	// The output will be the volume name if it exists, or empty if it doesn't.
	// We trim space and check if the output matches the name.
	exists := strings.TrimSpace(output) == name
	return exists, nil
}

func CreateVolume(name, size, baseDir string, labels map[string]string) (string, error) {
	// Check if the volume already exists
	exists, err := volumeExists(name)
	if err != nil {
		return "", fmt.Errorf("failed to check for existing volume: %v", err)
	}
	if exists {
		return "", fmt.Errorf("volume '%s' already exists", name)
	}

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

	lostAndFoundPath := filepath.Join(dataPath, "lost+found")
	log.Printf("Removing lost+found directory: %s", lostAndFoundPath)
	if err := runCommand("sudo", "rm", "-rf", lostAndFoundPath); err != nil {
		// Log as a warning instead of returning an error
		log.Printf("warning: failed to remove lost+found: %v", err)
	}

	log.Printf("Setting permissions for data directory: %s to 777", absDataPath)
	if err := runCommand("sudo", "chmod", "-R", "777", absDataPath); err != nil {
		return "", fmt.Errorf("chmod failed: %v", err)
	}

	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		log.Printf("Setting ownership for data directory: %s to %s", dataPath, sudoUser)
		if err := runCommand("sudo", "chown", "-R", fmt.Sprintf("%s:%s", sudoUser, sudoUser), dataPath); err != nil {
			return "", fmt.Errorf("chown failed: %v", err)
		}
	} // for dev
	log.Printf("Registering docker volume: %s", name)

	dockerArgs := []string{
		"docker", "volume", "create",
		"--name", name,
		"--opt", fmt.Sprintf("device=%s", absDataPath),
		"--opt", "type=none",
		"--opt", "o=bind",
	}

	for key, value := range labels {
		dockerArgs = append(dockerArgs, "--label", fmt.Sprintf("%s=%s", key, value))
	}

	if err := runCommand(dockerArgs[0], dockerArgs[1:]...); err != nil {
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

func GetVolumeStats(name, baseDir string) (*VolumeStats, error) {
	volumePath := filepath.Join(baseDir, name)
	dataPath := filepath.Join(volumePath, "_data")

	output, err := runCommandWithOutput("df", "-h", dataPath)
	if err != nil {
		return nil, fmt.Errorf("df command failed: %v", err)
	}

	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("invalid df output")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 6 {
		return nil, fmt.Errorf("invalid df output fields")
	}

	stats := &VolumeStats{
		Name:      name,
		Size:      formatSize(fields[1]),
		Used:      formatSize(fields[2]),
		Available: formatSize(fields[3]),
		Usage:     fields[4],
		MountPath: fields[5],
	}

	return stats, nil
}

func formatSize(size string) string {
	if len(size) < 1 {
		return size
	}
	lastChar := size[len(size)-1]
	if (lastChar >= 'A' && lastChar <= 'Z') || (lastChar >= 'a' && lastChar <= 'z') {
		value := size[:len(size)-1]
		return value + " " + string(lastChar) + "B"
	}
	return size
}

func GetAllVolumes(baseDir string) ([]*VolumeStats, error) {
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read base directory: %v", err)
	}

	var volumes []*VolumeStats
	for _, file := range files {
		if file.IsDir() {
			stats, err := GetVolumeStats(file.Name(), baseDir)
			if err != nil {
				// Log the error but continue, as some directories might not be volumes
				log.Printf("failed to get stats for %s: %v", file.Name(), err)
				continue
			}
			volumes = append(volumes, stats)
		}
	}

	return volumes, nil
}
