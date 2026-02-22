package wiki

import (
	"fmt"
	"strings"

	"github.com/GoooIce/repowiki/internal/config"
)

func BuildFullGeneratePrompt(cfg *config.Config) string {
	return fmt.Sprintf(`You are a technical documentation specialist. Generate a comprehensive repository wiki for this project.

OUTPUT REQUIREMENTS:
- Create documentation files in %s/%s/content/ directory
- Create a metadata file at %s/%s/meta/repowiki-metadata.json
- Each markdown file must follow this structure:
  1. Title as H1 heading
  2. <cite> block listing referenced source files with format: [filename](file://path/to/file)
  3. Table of Contents with anchor links
  4. Detailed content with code examples from the actual source
  5. Mermaid diagrams for architecture where appropriate

WIKI STRUCTURE — create these files/directories:
- System Overview.md — project purpose, high-level architecture
- Technology Stack.md — languages, frameworks, key dependencies
- Getting Started.md — setup, installation, running
- Backend Architecture/ — server structure, API design, database, etc.
- Frontend Architecture/ — UI components, state management, etc.
- Core Features/ — each major feature documented individually
- API Reference/ — endpoints, request/response formats
- Configuration Management.md — environment variables, config files

METADATA FORMAT for repowiki-metadata.json:
{
  "code_snippets": [
    {
      "id": "<md5 hash>",
      "path": "relative/path/to/file",
      "line_range": "1-100",
      "gmt_create": "<ISO 8601 timestamp>",
      "gmt_modified": "<ISO 8601 timestamp>"
    }
  ]
}

Analyze ALL source files. Be thorough. Include actual code references.
Do NOT modify any source code. Only create/modify files within %s/.`, cfg.WikiPath, cfg.Language, cfg.WikiPath, cfg.Language, cfg.WikiPath)
}

func BuildIncrementalPrompt(cfg *config.Config, changedFiles []string, affectedSections []string) string {
	fileList := "  - " + strings.Join(changedFiles, "\n  - ")

	sectionHint := ""
	if len(affectedSections) > 0 {
		sectionHint = fmt.Sprintf(`
POTENTIALLY AFFECTED WIKI SECTIONS (check and update these first):
  - %s
`, strings.Join(affectedSections, "\n  - "))
	}

	return fmt.Sprintf(`You are a technical documentation specialist. Update the repository wiki to reflect recent code changes.

CHANGED SOURCE FILES:
%s
%s
INSTRUCTIONS:
1. Read each changed source file to understand what was modified
2. Read the existing wiki pages in %s/%s/content/
3. Update ONLY the wiki sections affected by the code changes
4. If a changed file introduces new functionality not covered by existing pages, create a new page
5. Update %s/%s/meta/repowiki-metadata.json with any new or modified code snippet references
6. Preserve existing formatting: <cite> blocks, Table of Contents, mermaid diagrams
7. Do NOT modify any source code. Only modify files within %s/

Keep documentation accurate and synchronized with the current codebase.`, fileList, sectionHint, cfg.WikiPath, cfg.Language, cfg.WikiPath, cfg.Language, cfg.WikiPath)
}
