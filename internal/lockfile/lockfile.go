package lockfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const lockFileName = ".repowiki.lock"

func lockPath(gitRoot string) string {
	return filepath.Join(gitRoot, ".repowiki", lockFileName)
}

func Acquire(gitRoot string) error {
	lp := lockPath(gitRoot)

	if err := os.MkdirAll(filepath.Dir(lp), 0755); err != nil {
		return fmt.Errorf("failed to create lock dir: %w", err)
	}

	// Check for stale lock first
	if IsLocked(gitRoot) {
		if isStale(lp) {
			os.Remove(lp)
		} else {
			return fmt.Errorf("another repowiki process is running (lock: %s)", lp)
		}
	}

	f, err := os.OpenFile(lp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("another repowiki process is running")
		}
		return fmt.Errorf("failed to create lock: %w", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "%d\n%s\n", os.Getpid(), time.Now().UTC().Format(time.RFC3339))
	return nil
}

func Release(gitRoot string) {
	os.Remove(lockPath(gitRoot))
}

func IsLocked(gitRoot string) bool {
	_, err := os.Stat(lockPath(gitRoot))
	return err == nil
}

func isStale(lp string) bool {
	data, err := os.ReadFile(lp)
	if err != nil {
		return true
	}

	lines := strings.SplitN(string(data), "\n", 3)
	if len(lines) < 1 {
		return true
	}

	pid, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return true
	}

	// Check if process is still running
	proc, err := os.FindProcess(pid)
	if err != nil {
		return true
	}

	// On Unix, FindProcess always succeeds. Send signal 0 to check.
	err = proc.Signal(os.Signal(nil))
	if err != nil {
		return true // Process not running
	}

	// Check age - stale if older than 30 minutes
	if len(lines) >= 2 {
		ts, err := time.Parse(time.RFC3339, strings.TrimSpace(lines[1]))
		if err == nil && time.Since(ts) > 30*time.Minute {
			return true
		}
	}

	return false
}
