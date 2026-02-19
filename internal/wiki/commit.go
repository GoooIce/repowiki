package wiki

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ikrasnodymov/repowiki/internal/config"
	"github.com/ikrasnodymov/repowiki/internal/git"
)

const sentinelFile = ".committing"

func sentinelPath(gitRoot string) string {
	return filepath.Join(config.Dir(gitRoot), sentinelFile)
}

// IsSentinelPresent checks if a wiki commit is in progress (loop prevention).
func IsSentinelPresent(gitRoot string) bool {
	_, err := os.Stat(sentinelPath(gitRoot))
	return err == nil
}

// CommitChanges stages and commits wiki changes with loop prevention.
func CommitChanges(gitRoot string, cfg *config.Config, description string) error {
	wikiDir := filepath.Join(gitRoot, cfg.WikiPath)

	// Check if there are any changes to commit
	hasChanges, err := git.HasChanges(gitRoot, wikiDir)
	if err != nil || !hasChanges {
		return nil // Nothing to commit
	}

	// Write sentinel file (loop prevention layer 1)
	sp := sentinelPath(gitRoot)
	if err := os.WriteFile(sp, []byte(strconv.Itoa(os.Getpid())), 0644); err != nil {
		return fmt.Errorf("failed to write sentinel: %w", err)
	}
	defer os.Remove(sp)

	// Stage wiki files
	if err := git.StageFiles(gitRoot, []string{wikiDir}); err != nil {
		return fmt.Errorf("failed to stage wiki files: %w", err)
	}

	// Also stage config (updated last_run, last_commit_hash)
	configPath := config.Path(gitRoot)
	if _, err := os.Stat(configPath); err == nil {
		git.StageFiles(gitRoot, []string{configPath})
	}

	// Commit with recognizable prefix
	message := fmt.Sprintf("%s %s", cfg.CommitPrefix, description)
	if err := git.Commit(gitRoot, message); err != nil {
		return fmt.Errorf("failed to commit wiki: %w", err)
	}

	return nil
}
