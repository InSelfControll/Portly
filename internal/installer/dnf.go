package installer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func (i *Installer) installWithDNF(ctx context.Context, config ProductConfig) error {
	for _, dep := range config.Dependencies {
		cmd := exec.CommandContext(ctx, "dnf", "install", "-y", dep)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: failed to install %s: %v\n", dep, err)
		}
	}

	cmd := exec.CommandContext(ctx, "dnf", "install", "-y", config.PackageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return i.runPostInstall(ctx, config)
}

func (i *Installer) installWithYUM(ctx context.Context, config ProductConfig) error {
	for _, dep := range config.Dependencies {
		cmd := exec.CommandContext(ctx, "yum", "install", "-y", dep)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: failed to install %s: %v\n", dep, err)
		}
	}

	cmd := exec.CommandContext(ctx, "yum", "install", "-y", config.PackageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return i.runPostInstall(ctx, config)
}
