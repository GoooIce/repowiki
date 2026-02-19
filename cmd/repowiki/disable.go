package main

import (
	"fmt"
	"os"

	"github.com/ikrasnodymov/repowiki/internal/config"
	"github.com/ikrasnodymov/repowiki/internal/git"
	"github.com/ikrasnodymov/repowiki/internal/hook"
)

func handleDisable(args []string) {
	gitRoot, err := git.FindRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: not a git repository\n")
		os.Exit(1)
	}

	// Remove hook
	if err := hook.Uninstall(gitRoot); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing hook: %v\n", err)
		os.Exit(1)
	}

	// Update config
	cfg, err := config.Load(gitRoot)
	if err == nil {
		cfg.Enabled = false
		config.Save(gitRoot, cfg)
	}

	fmt.Printf("repowiki disabled in %s\n", gitRoot)
	fmt.Printf("Wiki files in .qoder/repowiki/ are preserved.\n")
}
