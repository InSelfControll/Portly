package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	// Main Menu Items
	items := []list.Item{
		menuItem{"Add Rule Setup", "Configure new firewall rules", ScreenAddRuleSelect},
		menuItem{"List Rules", "View and manage existing rules", ScreenListRules},
		menuItem{"Firewall", "Start/stop or install firewall", ScreenFirewall},
		menuItem{"Security", "Manage SELinux/AppArmor", ScreenSecurity},
		menuItem{"System Status", "View system and provider status", ScreenStatus},
		menuItem{"Check Configuration", "Verify system configuration", ScreenCheck},
		menuItem{"Quit", "Exit Portly", -1},
	}

	subItems := []list.Item{
		menuItem{"NAT Rule", "Forward external port to internal IP/Port", ScreenAddNATRule},
		menuItem{"Open Port", "Open a port for all incoming traffic", ScreenOpenPort},
		menuItem{"IP Restricted Port", "Open a specific port only for a specific IP", ScreenOpenIPPort},
		menuItem{"Open All Ports for IP", "Allow all traffic from a specific IP", ScreenOpenIP},
	}

	menuList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	menuList.Title = "Portly - Main Menu"
	menuList.SetShowStatusBar(false)
	menuList.SetFilteringEnabled(false)
	menuList.Styles.Title = styles.Title

	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = lipgloss.NewStyle().Foreground(lipgloss.Color(styles.TextColor)).Padding(0, 0, 0, 2)
	delegate.Styles.NormalDesc = lipgloss.NewStyle().Foreground(lipgloss.Color(styles.MutedColor)).Padding(0, 0, 0, 2)
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(styles.PrimaryColor)).Padding(0, 0, 0, 2).Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color(styles.PrimaryColor))
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().Foreground(lipgloss.Color(styles.SecondaryColor)).Padding(0, 0, 0, 2).Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color(styles.PrimaryColor))

	subMenuList := list.New(subItems, delegate, 20, 10)
	subMenuList.Title = "Add Rule - Select Type"
	subMenuList.SetShowStatusBar(false)
	subMenuList.SetFilteringEnabled(false)
	subMenuList.Styles.Title = styles.Title

	return &Model{
		ctx:             ctx,
		osInfo:          osInfo,
		provider:        provider,
		stateMgr:        stateMgr,
		screen:          ScreenMenu,
		menuItems:       items,
		menuList:        menuList,
		subItems:        subItems,
		ruleSubMenuList: subMenuList,
		addRuleForm:     NewAddRuleForm(),
		ruleViewMode:    "nat",
	}, nil
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return m.addRuleForm.Init()
}
