package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ikrasnodymov/repowiki/internal/config"
	"github.com/ikrasnodymov/repowiki/internal/git"
	"github.com/ikrasnodymov/repowiki/internal/wiki"
)

func handleUpdate(args []string) {
	fs := flag.NewFlagSet("update", flag.ExitOnError)
	commitHash := fs.String("commit", "", "specific commit hash to process")
	fromHook := fs.Bool("from-hook", false, "internal: hook-triggered run")
	fs.Parse(args)

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

	// Determine which commit to process
	hash := *commitHash
	if hash == "" {
		hash, err = git.HeadCommit(gitRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting HEAD: %v\n", err)
			os.Exit(1)
		}
	}

	// Get changed files
	var changedFiles []string
	if cfg.LastCommitHash != "" && cfg.LastCommitHash != hash {
		changedFiles, err = git.ChangedFilesSince(gitRoot, cfg.LastCommitHash)
	} else {
		changedFiles, err = git.ChangedFilesInCommit(gitRoot, hash)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting changes: %v\n", err)
		os.Exit(1)
	}

	// Filter excluded paths
	changedFiles = filterExcluded(changedFiles, cfg.ExcludedPaths)

	if len(changedFiles) == 0 {
		if !*fromHook {
			fmt.Println("No relevant file changes detected.")
		}
		return
	}

	// Decide: full generate or incremental
	if !wiki.Exists(gitRoot, cfg) || len(changedFiles) > cfg.FullGenerateThreshold {
		if !*fromHook {
			fmt.Printf("Running full wiki generation (%d files changed)...\n", len(changedFiles))
		}
		if err := wiki.FullGenerate(gitRoot, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		if !*fromHook {
			fmt.Printf("Updating wiki for %d changed files...\n", len(changedFiles))
		}
		if err := wiki.IncrementalUpdate(gitRoot, cfg, changedFiles); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	// Update last run
	config.UpdateLastRun(gitRoot, hash)

	if !*fromHook {
		fmt.Println("Wiki update complete.")
	}
}

func filterExcluded(files []string, excluded []string) []string {
	var result []string
	for _, f := range files {
		skip := false
		for _, ex := range excluded {
			if len(f) >= len(ex) && f[:len(ex)] == ex {
				skip = true
				break
			}
		}
		if !skip {
			result = append(result, f)
		}
	}
	return result
}
