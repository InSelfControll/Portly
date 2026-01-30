package state

import (
	"fmt"
	"os"
	"time"

	"github.com/orchestrator/unified-firewall/internal/platform"
	"github.com/orchestrator/unified-firewall/pkg/models"
)

// NewManager creates a new state manager
func NewManager() (*Manager, error) {
	if err := platform.EnsureStateDir(); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %w", err)
	}

	m := &Manager{
		statePath: platform.GetStateFilePath(),
		state:     nil,
	}

	storage := NewFileStorage(m.statePath)
	state, err := storage.Load()
	if err != nil {
		if os.IsNotExist(err) {
			if err := m.Initialize(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		m.state = state
	}

	return m, nil
}

// Initialize creates a new state file
func (m *Manager) Initialize() error {
	osInfo, err := platform.DetectOS()
	if err != nil {
		return fmt.Errorf("failed to detect OS: %w", err)
	}

	m.state = &models.State{
		Version: "1.0.0",
		OS: models.OSInfo{
			Family:       string(osInfo.Family),
			Distribution: osInfo.Distribution,
			Version:      osInfo.Version,
			Codename:     osInfo.Codename,
		},
		Rules:       []models.AppliedRule{},
		Products:    make(map[string]models.ProductInfo),
		LastUpdated: time.Now().UTC().Format(time.RFC3339),
	}

	return m.Save()
}

// Save writes state to disk
func (m *Manager) Save() error {
	if m.state == nil {
		return fmt.Errorf("state is not initialized")
	}

	m.state.LastUpdated = time.Now().UTC().Format(time.RFC3339)
	storage := NewFileStorage(m.statePath)
	return storage.Save(m.state)
}

// GetState returns the current state
func (m *Manager) GetState() *models.State {
	return m.state
}
