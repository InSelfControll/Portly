package installer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func (i *Installer) installWithBrew(ctx context.Context, config ProductConfig) error {
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("Homebrew not installed: https://brew.sh")
	}

	fmt.Println("Updating Homebrew...")
	cmd := exec.CommandContext(ctx, "brew", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: brew update failed: %v\n", err)
	}

	formula := config.BrewFormula
	if formula == "" {
		formula = config.PackageName
	}

	cmd = exec.CommandContext(ctx, "brew", "install", formula)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	for _, dep := range config.Dependencies {
		cmd := exec.CommandContext(ctx, "brew", "install", dep)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: failed to install %s: %v\n", dep, err)
		}
	}

	return i.runPostInstall(ctx, config)
}
