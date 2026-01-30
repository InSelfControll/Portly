package drivers

import (
	"context"

	"github.com/orchestrator/unified-firewall/internal/drivers/firewalld"
	"github.com/orchestrator/unified-firewall/internal/drivers/nftables"
	"github.com/orchestrator/unified-firewall/internal/drivers/pf"
	"github.com/orchestrator/unified-firewall/pkg/models"
)

// Provider defines the interface that all OS-specific drivers must implement
type Provider interface {
	Name() string
	IsAvailable() bool

	// Product installation
	IsProductInstalled(ctx context.Context, name string) (models.ProductInfo, error)
	GetInstalledProducts(ctx context.Context) ([]models.ProductInfo, error)

	// NAT rules
	ApplyNAT(ctx context.Context, rule models.NATRule) error
	RemoveNAT(ctx context.Context, ruleID string) error
	ListNATRules(ctx context.Context) ([]models.NATRule, error)
	CheckConflicts(ctx context.Context, port int, proto models.Protocol) error

	// Firewall rules (port opening)
	OpenPort(ctx context.Context, rule models.FirewallRule) error
	OpenPortForIP(ctx context.Context, rule models.FirewallRule) error
	TrustIP(ctx context.Context, rule models.FirewallRule) error
	ClosePort(ctx context.Context, ruleID string) error
	ListFirewallRules(ctx context.Context) ([]models.FirewallRule, error)

	// Security policies
	EnsureSecurityPolicy(ctx context.Context, product string, policy models.SecurityPolicy) error
	RemoveSecurityPolicy(ctx context.Context, product string) error
}

// ProviderFactory creates providers based on OS detection
type ProviderFactory struct {
	providers []Provider
}

// NewProviderFactory creates a new factory with all available providers
func NewProviderFactory() *ProviderFactory {
	factory := &ProviderFactory{
		providers: make([]Provider, 0),
	}

	candidates := []Provider{
		firewalld.New(),
		nftables.New(),
		pf.New(),
	}

	for _, p := range candidates {
		if p != nil && p.IsAvailable() {
			factory.providers = append(factory.providers, p)
		}
	}

	return factory
}

// GetProvider returns the appropriate provider for the current system
func (f *ProviderFactory) GetProvider() (Provider, error) {
	if len(f.providers) == 0 {
		return nil, ErrNoProviderAvailable
	}

	return f.providers[0], nil
}

// GetAllProviders returns all available providers
func (f *ProviderFactory) GetAllProviders() []Provider {
	return f.providers
}
