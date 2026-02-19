package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/ikrasnodymov/repowiki/internal/config"
	"github.com/ikrasnodymov/repowiki/internal/git"
	"github.com/ikrasnodymov/repowiki/internal/lockfile"
	"github.com/ikrasnodymov/repowiki/internal/wiki"
)

// handleHooks is the entry point called by the git post-commit hook.
// It runs loop prevention checks and spawns a background update process.
func handleHooks(args []string) {
	if len(args) == 0 || args[0] != "post-commit" {
		return
	}

	gitRoot, err := git.FindRoot()
	if err != nil {
		return
	}

	// Loop prevention layer 1: sentinel file
	if wiki.IsSentinelPresent(gitRoot) {
		return
	}

	// Loop prevention layer 2: lock file
	if lockfile.IsLocked(gitRoot) {
		return
	}

	// Load config
	cfg, err := config.Load(gitRoot)
	if err != nil || !cfg.Enabled {
		return
	}

	// Get current commit
	commitHash, err := git.HeadCommit(gitRoot)
	if err != nil {
		return
	}

	// Loop prevention layer 3: check commit message prefix
	commitMsg, err := git.CommitMessage(gitRoot, commitHash)
	if err != nil {
		return
	}
	if strings.HasPrefix(strings.TrimSpace(commitMsg), cfg.CommitPrefix) {
		return
	}

	// All checks passed — spawn background update process
	spawnBackground(gitRoot, commitHash)
}

// spawnBackground launches `repowiki update --from-hook --commit <hash>` as a
// detached process so the user's terminal is not blocked.
func spawnBackground(gitRoot string, commitHash string) {
	self, err := os.Executable()
	if err != nil {
		return
	}

	logDir := config.LogPath(gitRoot)
	os.MkdirAll(logDir, 0755)

	logFile, err := os.OpenFile(
		fmt.Sprintf("%s/hook.log", logDir),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return
	}

	cmd := exec.Command(self, "update", "--from-hook", "--commit", commitHash)
	cmd.Dir = gitRoot
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	cmd.Start()
	// Do NOT call cmd.Wait() — let it run independently
	logFile.Close()
}
