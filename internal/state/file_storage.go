package state

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// FileStorage implements state storage using a JSON file
type FileStorage struct {
	path string
}

// NewFileStorage creates a new file-based storage
func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path}
}

// Load reads state from disk
func (fs *FileStorage) Load() (*models.State, error) {
	data, err := os.ReadFile(fs.path)
	if err != nil {
		return nil, err
	}

	var state models.State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return &state, nil
}

// Save writes state to disk atomically
func (fs *FileStorage) Save(state *models.State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	tempPath := fs.path + ".tmp"
	if err := os.WriteFile(tempPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return os.Rename(tempPath, fs.path)
}
