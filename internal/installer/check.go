package installer

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/orchestrator/unified-firewall/internal/platform"
	"github.com/orchestrator/unified-firewall/pkg/models"
)

// NewInstaller creates a new installer
func NewInstaller(autoAccept bool) (*Installer, error) {
	osInfo, err := platform.DetectOS()
	if err != nil {
		return nil, fmt.Errorf("failed to detect OS: %w", err)
	}

	return &Installer{
		osInfo:     osInfo,
		autoAccept: autoAccept,
	}, nil
}

// CheckProduct checks if a product is installed and prompts for installation
func (i *Installer) CheckProduct(ctx context.Context, productName string) (models.ProductInfo, error) {
	config, supported := SupportedProducts[productName]
	if !supported {
		return i.checkBinary(ctx, productName)
	}

	path, exists := platform.IsProductInstalled(productName)
	if exists {
		return models.ProductInfo{
			Name:        productName,
			Path:        path,
			IsInstalled: true,
		}, nil
	}

	if !i.autoAccept {
		shouldInstall, err := i.promptUser(config)
		if err != nil {
			return models.ProductInfo{}, err
		}
		if !shouldInstall {
			return models.ProductInfo{
				Name:        productName,
				IsInstalled: false,
			}, fmt.Errorf("user declined installation of %s", productName)
		}
	}

	fmt.Printf("Installing %s...\n", config.DisplayName)
	if err := i.InstallProduct(ctx, config); err != nil {
		return models.ProductInfo{}, fmt.Errorf("failed to install %s: %w", productName, err)
	}

	path, exists = platform.IsProductInstalled(productName)
	if !exists {
		return models.ProductInfo{}, fmt.Errorf("installation completed but %s not found", productName)
	}

	fmt.Printf("✓ %s installed at %s\n", config.DisplayName, path)

	return models.ProductInfo{
		Name:        productName,
		Path:        path,
		IsInstalled: true,
	}, nil
}

func (i *Installer) promptUser(config ProductConfig) (bool, error) {
	pm := i.osInfo.PackageManager()

	fmt.Printf("\n⚠️  '%s' is not installed.\n", config.DisplayName)
	fmt.Printf("Install using %s? (y/n): ", pm)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes", nil
}

func (i *Installer) checkBinary(ctx context.Context, name string) (models.ProductInfo, error) {
	path, exists := platform.IsProductInstalled(name)

	info := models.ProductInfo{
		Name:        name,
		Path:        path,
		IsInstalled: exists,
	}

	if exists {
		cmd := exec.CommandContext(ctx, name, "--version")
		output, err := cmd.Output()
		if err == nil {
			info.Version = strings.TrimSpace(string(output))
		}
	}

	return info, nil
}
