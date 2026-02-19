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

	hash := *commitHash
	if hash == "" {
		hash, err = git.HeadCommit(gitRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting HEAD: %v\n", err)
			os.Exit(1)
		}
	}

	if err := runUpdateCycle(gitRoot, cfg, hash, *fromHook); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// When running from hook, check if new commits arrived during generation.
	// If a commit happened while we held the lock, its hook exited silently.
	// Re-run to pick up those missed changes.
	if *fromHook {
		for i := 0; i < 5; i++ { // cap retries to avoid runaway loops
			cfg, err = config.Load(gitRoot)
			if err != nil {
				break
			}
			head, err := git.HeadCommit(gitRoot)
			if err != nil {
				break
			}
			if !hasUnprocessedCommits(gitRoot, cfg, head) {
				break
			}
			if err := runUpdateCycle(gitRoot, cfg, head, true); err != nil {
				break
			}
		}
	}

	if !*fromHook {
		fmt.Println("Wiki update complete.")
	}
}

// hasUnprocessedCommits checks if there are non-repowiki commits after the
// last processed commit.
func hasUnprocessedCommits(gitRoot string, cfg *config.Config, head string) bool {
	if cfg.LastCommitHash == "" || cfg.LastCommitHash == head {
		return false
	}
	// Check that the gap contains actual code changes, not just repowiki commits
	files, err := git.ChangedFilesSince(gitRoot, cfg.LastCommitHash)
	if err != nil {
		return false
	}
	files = filterExcluded(files, cfg.ExcludedPaths)
	return len(files) > 0
}

// runUpdateCycle performs a single update cycle: detect changes, run generation.
func runUpdateCycle(gitRoot string, cfg *config.Config, hash string, fromHook bool) error {
	var changedFiles []string
	var err error
	if cfg.LastCommitHash != "" && cfg.LastCommitHash != hash {
		changedFiles, err = git.ChangedFilesSince(gitRoot, cfg.LastCommitHash)
	} else {
		changedFiles, err = git.ChangedFilesInCommit(gitRoot, hash)
	}
	if err != nil {
		return fmt.Errorf("detecting changes: %w", err)
	}

	changedFiles = filterExcluded(changedFiles, cfg.ExcludedPaths)

	if len(changedFiles) == 0 {
		if !fromHook {
			fmt.Println("No relevant file changes detected.")
		}
		return nil
	}

	if !wiki.Exists(gitRoot, cfg) || len(changedFiles) > cfg.FullGenerateThreshold {
		if !fromHook {
			fmt.Printf("Running full wiki generation (%d files changed)...\n", len(changedFiles))
		}
		return wiki.FullGenerate(gitRoot, cfg, hash)
	}

	if !fromHook {
		fmt.Printf("Updating wiki for %d changed files...\n", len(changedFiles))
	}
	return wiki.IncrementalUpdate(gitRoot, cfg, changedFiles, hash)
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
