package platform

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// DetectOS detects the current operating system
func DetectOS() (*OSInfo, error) {
	switch runtime.GOOS {
	case "linux":
		return detectLinux()
	case "darwin":
		return detectDarwin()
	default:
		return &OSInfo{
			Family:       FamilyUnknown,
			Distribution: runtime.GOOS,
			Version:      "unknown",
		}, nil
	}
}

// detectLinux parses /etc/os-release to determine the Linux distribution
func detectLinux() (*OSInfo, error) {
	info := &OSInfo{}

	file, err := os.Open("/etc/os-release")
	if err != nil {
		file, err = os.Open("/usr/lib/os-release")
		if err != nil {
			return nil, fmt.Errorf("unable to detect OS: %w", err)
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := strings.Trim(parts[1], `"`)

		switch key {
		case "ID":
			info.ID = value
		case "ID_LIKE":
			info.IDLike = value
		case "VERSION_ID":
			info.Version = value
		case "VERSION_CODENAME":
			info.Codename = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading os-release: %w", err)
	}

	info.Distribution = info.ID
	info.Family = determineFamily(info.ID, info.IDLike)

	return info, nil
}

// detectDarwin detects macOS version
func detectDarwin() (*OSInfo, error) {
	info := &OSInfo{
		Family:       FamilyDarwin,
		Distribution: "macos",
	}

	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err == nil {
		info.Version = strings.TrimSpace(string(output))
	}

	return info, nil
}
