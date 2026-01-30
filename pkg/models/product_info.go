package models

// ProductInfo contains information about an installed product
type ProductInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Version     string `json:"version"`
	IsInstalled bool   `json:"is_installed"`
}
