package wiki

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/ikrasnodymov/repowiki/internal/config"
)

func FindQoderCLI(cfg *config.Config) (string, error) {
	// 1. Use config override
	if cfg.QoderCLIPath != "" && cfg.QoderCLIPath != "qodercli" {
		if _, err := os.Stat(cfg.QoderCLIPath); err == nil {
			return cfg.QoderCLIPath, nil
		}
	}

	// 2. Check PATH
	if path, err := exec.LookPath("qodercli"); err == nil {
		return path, nil
	}

	// 3. Check known macOS locations
	if runtime.GOOS == "darwin" {
		knownPaths := []string{
			"/Applications/Qoder.app/Contents/Resources/app/resources/bin/aarch64_darwin/qodercli",
			"/Applications/Qoder.app/Contents/Resources/app/resources/bin/x86_64_darwin/qodercli",
		}
		for _, p := range knownPaths {
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}
		}
	}

	// 4. Check known Linux locations
	if runtime.GOOS == "linux" {
		knownPaths := []string{
			"/usr/bin/qodercli",
			"/usr/local/bin/qodercli",
		}
		for _, p := range knownPaths {
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}
		}
	}

	return "", fmt.Errorf("qodercli not found; install Qoder or set qodercli_path in .repowiki/config.json")
}

func runQoder(cfg *config.Config, gitRoot string, prompt string) (string, error) {
	cliPath, err := FindQoderCLI(cfg)
	if err != nil {
		return "", err
	}

	args := []string{
		"-p", prompt,
		"-q",
		"-w", gitRoot,
		"--max-turns", strconv.Itoa(cfg.MaxTurns),
		"--dangerously-skip-permissions",
		"--allowed-tools", "Read,Write,Edit,Glob,Grep,Bash",
	}

	if cfg.Model != "" && cfg.Model != "auto" {
		args = append(args, "--model", cfg.Model)
	}

	cmd := exec.Command(cliPath, args...)
	cmd.Dir = gitRoot

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("qodercli error: %w\nstderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
