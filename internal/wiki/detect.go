package wiki

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoooIce/repowiki/internal/config"
)

type codeSnippet struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	LineRange string `json:"line_range"`
}

type metadata struct {
	CodeSnippets []codeSnippet `json:"code_snippets"`
}

// AffectedSections determines which wiki sections need updating based on changed files.
// It uses the metadata reverse index and heuristic path matching.
func AffectedSections(gitRoot string, cfg *config.Config, changedFiles []string) []string {
	affected := map[string]bool{}

	// 1. Build reverse index from metadata
	reverseIdx := buildReverseIndex(gitRoot, cfg)
	for _, f := range changedFiles {
		if pages, ok := reverseIdx[f]; ok {
			for _, p := range pages {
				affected[p] = true
			}
		}
	}

	// 2. Heuristic path matching
	for _, f := range changedFiles {
		for _, section := range heuristicMatch(f) {
			affected[section] = true
		}
	}

	result := make([]string, 0, len(affected))
	for s := range affected {
		result = append(result, s)
	}
	return result
}

// buildReverseIndex reads repowiki-metadata.json and cross-references with
// wiki content files to map source files -> wiki pages that reference them.
func buildReverseIndex(gitRoot string, cfg *config.Config) map[string][]string {
	idx := map[string][]string{}

	metaPath := filepath.Join(gitRoot, cfg.WikiPath, cfg.Language, "meta", "repowiki-metadata.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return idx
	}

	var meta metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return idx
	}

	// Collect all source file paths from metadata
	sourceFiles := map[string]bool{}
	for _, s := range meta.CodeSnippets {
		sourceFiles[s.Path] = true
	}

	// Scan wiki content files for <cite> blocks referencing source files
	contentDir := filepath.Join(gitRoot, cfg.WikiPath, cfg.Language, "content")
	scanWikiContent(contentDir, "", sourceFiles, idx)

	return idx
}

func scanWikiContent(dir string, relDir string, sourceFiles map[string]bool, idx map[string][]string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if e.IsDir() {
			subRel := filepath.Join(relDir, e.Name())
			scanWikiContent(filepath.Join(dir, e.Name()), subRel, sourceFiles, idx)
			continue
		}

		if !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		wikiPage := filepath.Join(relDir, e.Name())
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}

		content := string(data)
		// Look for file:// references in <cite> blocks
		for srcFile := range sourceFiles {
			if strings.Contains(content, "file://"+srcFile) || strings.Contains(content, srcFile) {
				idx[srcFile] = append(idx[srcFile], wikiPage)
			}
		}
	}
}

func heuristicMatch(filePath string) []string {
	var sections []string
	lower := strings.ToLower(filePath)

	switch {
	case strings.Contains(lower, "backend/") || strings.Contains(lower, "server/") || strings.Contains(lower, "src/api/"):
		sections = append(sections, "Backend Architecture")
	case strings.Contains(lower, "frontend/") || strings.Contains(lower, "src/components/") || strings.Contains(lower, "src/app/"):
		sections = append(sections, "Frontend Architecture")
	}

	if strings.Contains(lower, "api/") || strings.Contains(lower, "routes/") || strings.Contains(lower, "endpoints/") {
		sections = append(sections, "API Reference")
	}

	if strings.Contains(lower, "config") || strings.Contains(lower, ".env") || strings.Contains(lower, "settings") {
		sections = append(sections, "Configuration Management")
	}

	if strings.HasSuffix(lower, "readme.md") || strings.HasSuffix(lower, "package.json") || strings.HasSuffix(lower, "pyproject.toml") {
		sections = append(sections, "System Overview")
	}

	if strings.Contains(lower, "auth") || strings.Contains(lower, "security") {
		sections = append(sections, "Authentication and Security")
	}

	if strings.Contains(lower, "database/") || strings.Contains(lower, "models/") || strings.Contains(lower, "migrations/") {
		sections = append(sections, "Backend Architecture")
	}

	return sections
}
