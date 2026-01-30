package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/orchestrator/unified-firewall/internal/drivers"
	"github.com/orchestrator/unified-firewall/internal/platform"
	"github.com/orchestrator/unified-firewall/internal/state"
	"github.com/orchestrator/unified-firewall/internal/tui/styles"
)

// New creates a new TUI model
func New(ctx context.Context) (*Model, error) {
	osInfo, err := platform.DetectOS()
	if err != nil {
		return nil, fmt.Errorf("failed to detect OS: %w", err)
	}

	factory := drivers.NewProviderFactory()
	provider, _ := factory.GetProvider()

	stateMgr, err := state.NewManager()
	if err != nil {
		stateMgr = nil
	}

	items := []list.Item{
		menuItem{"Add NAT Rule", "Create a new port forwarding rule", ScreenAddRule},
		menuItem{"List Rules", "View and manage existing rules", ScreenListRules},
		menuItem{"Firewall", "Start/stop or install firewall", ScreenFirewall},
		menuItem{"Security", "Manage SELinux/AppArmor", ScreenSecurity},
		menuItem{"System Status", "View system and provider status", ScreenStatus},
		menuItem{"Check Configuration", "Verify system configuration", ScreenCheck},
		menuItem{"Quit", "Exit Portly", -1},
	}

	menuList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	menuList.Title = "Portly - Main Menu"
	menuList.SetShowStatusBar(false)
	menuList.SetFilteringEnabled(false)
	menuList.Styles.Title = styles.Title

	return &Model{
		ctx:         ctx,
		osInfo:      osInfo,
		provider:    provider,
		stateMgr:    stateMgr,
		screen:      ScreenMenu,
		menuItems:   items,
		menuList:    menuList,
		addRuleForm: NewAddRuleForm(),
	}, nil
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return m.addRuleForm.Init()
}
