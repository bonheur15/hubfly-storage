package volume

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Helper to generate a unique volume name for tests
func generateVolumeName() string {
	return fmt.Sprintf("test-volume-%d", time.Now().UnixNano())
}

func TestCreateAndDeleteVolume(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "hubfly-volume-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	volumeName := generateVolumeName()
	volumeSize := "100M" // Use a small size for testing

	t.Run("Create Volume Successfully", func(t *testing.T) {
		createdName, err := CreateVolume(volumeName, volumeSize, tempDir)
		if err != nil {
			t.Fatalf("CreateVolume failed: %v", err)
		}
		if createdName != volumeName {
			t.Errorf("expected created volume name %q, got %q", volumeName, createdName)
		}

		// Verify volume exists in Docker
		exists, err := volumeExists(volumeName)
		if err != nil {
			t.Fatalf("volumeExists failed: %v", err)
		}
		if !exists {
			t.Error("volume does not exist after creation")
		}

		// Verify volume path exists
		volumePath := filepath.Join(tempDir, volumeName)
		if _, err := os.Stat(volumePath); os.IsNotExist(err) {
			t.Errorf("volume directory %q does not exist", volumePath)
		}
	})

	t.Run("Prevent Duplicate Volume Creation", func(t *testing.T) {
		_, err := CreateVolume(volumeName, volumeSize, tempDir)
		if err == nil {
			t.Error("expected error when creating duplicate volume, but got none")
		}
		expectedErrorMsg := fmt.Sprintf("volume '%s' already exists", volumeName)
		if err != nil && !strings.Contains(err.Error(), expectedErrorMsg) {
			t.Errorf("expected error message to contain %q, got %q", expectedErrorMsg, err.Error())
		}
	})

	t.Run("Delete Volume Successfully", func(t *testing.T) {
		err := DeleteVolume(volumeName, tempDir)
		if err != nil {
			t.Fatalf("DeleteVolume failed: %v", err)
		}

		// Verify volume no longer exists in Docker
		exists, err := volumeExists(volumeName)
		if err != nil {
			t.Fatalf("volumeExists failed after deletion: %v", err)
		}
		if exists {
			t.Error("volume still exists after deletion")
		}

		// Verify volume path is removed
		volumePath := filepath.Join(tempDir, volumeName)
		if _, err := os.Stat(volumePath); !os.IsNotExist(err) {
			t.Errorf("volume directory %q still exists after deletion", volumePath)
		}
	})
}

// Test for GetVolumeStats
func TestGetVolumeStats(t *testing.T) {
	// This test relies on a successfully created volume, similar to TestCreateAndDeleteVolume
	tempDir, err := ioutil.TempDir("", "hubfly-volume-test-stats")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	volumeName := generateVolumeName()
	volumeSize := "100M"

	_, err = CreateVolume(volumeName, volumeSize, tempDir)
	if err != nil {
		t.Fatalf("Failed to create volume for stats test: %v", err)
	}
	defer DeleteVolume(volumeName, tempDir) // Ensure cleanup even if stats test fails

	stats, err := GetVolumeStats(volumeName, tempDir)
	if err != nil {
		t.Fatalf("GetVolumeStats failed: %v", err)
	}

	if stats == nil {
		t.Fatal("GetVolumeStats returned nil stats")
	}

	if stats.Name != volumeName {
		t.Errorf("Expected stats name %q, got %q", volumeName, stats.Name)
	}
	// Further checks for Size, Used, Available, Usage, MountPath can be added
	// but are highly dependent on the `df -h` output format and system specifics.
	// For now, just ensuring the call succeeds and name matches is a good start.
	if !strings.Contains(stats.MountPath, filepath.Join(tempDir, volumeName, "_data")) {
		t.Errorf("MountPath %q does not contain expected volume path %q", stats.MountPath, filepath.Join(tempDir, volumeName, "_data"))
	}
}

// Test runCommand and runCommandWithOutput for basic functionality
func TestRunCommands(t *testing.T) {
	t.Run("runCommand success", func(t *testing.T) {
		err := runCommand("echo", "hello")
		if err != nil {
			t.Errorf("runCommand failed: %v", err)
		}
	})

	t.Run("runCommand failure", func(t *testing.T) {
		err := runCommand("false") // command that exits with non-zero status
		if err == nil {
			t.Error("runCommand did not return an error for a failing command")
		}
	})

	t.Run("runCommandWithOutput success", func(t *testing.T) {
		output, err := runCommandWithOutput("echo", "test output")
		if err != nil {
			t.Errorf("runCommandWithOutput failed: %v", err)
		}
		if strings.TrimSpace(output) != "test output" {
			t.Errorf("expected 'test output', got %q", output)
		}
	})

	t.Run("runCommandWithOutput failure", func(t *testing.T) {
		output, err := runCommandWithOutput("bash", "-c", "echo error_msg >&2; exit 1")
		if err == nil {
			t.Error("runCommandWithOutput did not return an error for a failing command")
		}
		if !strings.Contains(output, "error_msg") {
			t.Errorf("expected output to contain 'error_msg', got %q", output)
		}
	})
}