package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/orchestrator/unified-firewall/internal/drivers"
	"github.com/orchestrator/unified-firewall/internal/platform"
	"github.com/orchestrator/unified-firewall/internal/state"
	"github.com/orchestrator/unified-firewall/pkg/models"
)

// Screen represents the current UI screen
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenAddRule
	ScreenListRules
	ScreenStatus
	ScreenCheck
	ScreenLoading
	ScreenError
	ScreenSuccess
)

// Model is the main TUI model
type Model struct {
	ctx       context.Context
	provider  drivers.Provider
	stateMgr  *state.Manager
	osInfo    *platform.OSInfo
	
	screen        Screen
	width         int
	height        int
	lastError     error
	successMsg    string
	
	menuItems []list.Item
	menuList  list.Model
	addRuleForm   *AddRuleForm
	
	rules         []models.AppliedRule
	loadingMsg    string
}

type menuItem struct {
	title       string
	description string
	screen      Screen
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }

// Key bindings
var keys = struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
	Tab   key.Binding
}{
	Up:    key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "up")),
	Down:  key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "down")),
	Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Back:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Quit:  key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	Tab:   key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next")),
}
