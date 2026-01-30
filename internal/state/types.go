package state

import (
	"github.com/orchestrator/unified-firewall/pkg/models"
)

// Manager handles state persistence and retrieval
type Manager struct {
	statePath string
	state     *models.State
}

// StateStorage defines the interface for state storage operations
type StateStorage interface {
	Load() (*models.State, error)
	Save(state *models.State) error
}
