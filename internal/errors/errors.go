package errors

import "errors"

// Common errors
var (
	ErrNoProviderAvailable  = errors.New("no suitable provider available")
	ErrProductNotInstalled  = errors.New("product is not installed")
	ErrRuleExists           = errors.New("NAT rule already exists")
	ErrRuleNotFound         = errors.New("NAT rule not found")
	ErrPortConflict         = errors.New("port conflict detected")
	ErrPermissionDenied     = errors.New("root/admin privileges required")
	ErrSecurityPolicyFailed = errors.New("failed to apply security policy")
	ErrUnsupportedProduct   = errors.New("unsupported product")
)
