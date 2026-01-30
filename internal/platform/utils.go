package platform

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// CommandExists checks if a command exists in PATH
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// IsRoot checks if the current process is running as root
func IsRoot() bool {
	return os.Geteuid() == 0
}

// GetStateDir returns the directory for storing state files
func GetStateDir() string {
	if runtime.GOOS == "darwin" {
		return "/usr/local/var/lib/orchestrator"
	}
	return "/var/lib/orchestrator"
}

// EnsureStateDir creates the state directory if it doesn't exist
func EnsureStateDir() error {
	dir := GetStateDir()
	return os.MkdirAll(dir, 0755)
}

// GetStateFilePath returns the full path to the state file
func GetStateFilePath() string {
	return filepath.Join(GetStateDir(), "state.json")
}

// IsProductInstalled checks if a product binary exists in PATH
func IsProductInstalled(name string) (string, bool) {
	path, err := exec.LookPath(name)
	return path, err == nil
}
