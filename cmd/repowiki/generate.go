package main

import (
	"fmt"
	"os"

	"github.com/ikrasnodymov/repowiki/internal/config"
	"github.com/ikrasnodymov/repowiki/internal/git"
	"github.com/ikrasnodymov/repowiki/internal/wiki"
)

func handleGenerate(args []string) {
	gitRoot, err := git.FindRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: not a git repository\n")
		os.Exit(1)
	}

	cfg, err := config.Load(gitRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: repowiki not configured. Run 'repowiki enable' first.\n")
		os.Exit(1)
	}

	fmt.Println("Starting full wiki generation... (this may take several minutes)")

	if err := wiki.FullGenerate(gitRoot, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Update last run
	head, _ := git.HeadCommit(gitRoot)
	config.UpdateLastRun(gitRoot, head)

	fmt.Println("Wiki generation complete.")
}
