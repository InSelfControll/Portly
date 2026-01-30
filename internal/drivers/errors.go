// Package drivers implements OS-specific firewall drivers
package drivers

import (
	"github.com/orchestrator/unified-firewall/internal/errors"
)

// Re-export errors for backwards compatibility
var (
	ErrNoProviderAvailable  = errors.ErrNoProviderAvailable
	ErrProductNotInstalled  = errors.ErrProductNotInstalled
	ErrRuleExists           = errors.ErrRuleExists
	ErrRuleNotFound         = errors.ErrRuleNotFound
	ErrPortConflict         = errors.ErrPortConflict
	ErrPermissionDenied     = errors.ErrPermissionDenied
	ErrSecurityPolicyFailed = errors.ErrSecurityPolicyFailed
	ErrUnsupportedProduct   = errors.ErrUnsupportedProduct
)
