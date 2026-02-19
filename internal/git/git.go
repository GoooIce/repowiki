package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), string(exitErr.Stderr))
		}
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out)), nil
}

func FindRoot() (string, error) {
	return run("", "rev-parse", "--show-toplevel")
}

func FindRootFrom(dir string) (string, error) {
	return run(dir, "rev-parse", "--show-toplevel")
}

func HeadCommit(gitRoot string) (string, error) {
	return run(gitRoot, "rev-parse", "HEAD")
}

func CommitMessage(gitRoot string, hash string) (string, error) {
	return run(gitRoot, "log", "-1", "--pretty=%B", hash)
}

func ChangedFilesInCommit(gitRoot string, hash string) ([]string, error) {
	out, err := run(gitRoot, "diff-tree", "--no-commit-id", "--name-only", "-r", hash)
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

func ChangedFilesSince(gitRoot string, hash string) ([]string, error) {
	out, err := run(gitRoot, "diff", "--name-only", hash, "HEAD")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

func StageFiles(gitRoot string, paths []string) error {
	args := append([]string{"add"}, paths...)
	_, err := run(gitRoot, args...)
	return err
}

func Commit(gitRoot string, message string) error {
	_, err := run(gitRoot, "commit", "-m", message)
	return err
}

func HasChanges(gitRoot string, path string) (bool, error) {
	out, err := run(gitRoot, "status", "--porcelain", path)
	if err != nil {
		return false, err
	}
	return out != "", nil
}
