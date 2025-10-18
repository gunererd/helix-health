package main

import (
	"bufio"
	"strings"
)

// parseHealthOutput parses helix --health output into rows.
// Returns rows, column header, and config info lines.
func parseHealthOutput(output string) ([]Row, string, []string) {
	var rows []Row
	var currentRow *Row
	var header string
	var configLines []string
	var allLines []string

	// Collect all lines
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}

	if len(allLines) == 0 {
		return rows, header, configLines
	}

	// Phase 1: Collect config lines (non-empty lines from start)
	i := 0
	for i < len(allLines) && strings.TrimSpace(allLines[i]) != "" {
		configLines = append(configLines, allLines[i])
		i++
	}

	// Phase 2: Skip blank lines
	for i < len(allLines) && strings.TrimSpace(allLines[i]) == "" {
		i++
	}

	// Phase 3: Next non-empty line is the header
	if i < len(allLines) {
		header = allLines[i]
		i++
	}

	// Phase 4: Skip blank lines after header
	for i < len(allLines) && strings.TrimSpace(allLines[i]) == "" {
		i++
	}

	// Phase 5: Parse language rows
	for i < len(allLines) {
		line := allLines[i]

		// Check if this is a new row (line doesn't start with space or tab)
		if len(line) > 0 && line[0] != ' ' && line[0] != '\t' {
			// Save previous row
			if currentRow != nil {
				rows = append(rows, *currentRow)
			}

			// Start new row
			parts := strings.Fields(line)
			if len(parts) == 0 {
				i++
				continue
			}
			language := parts[0]

			tools := extractTools(line)
			searchTerms := []string{language}
			searchTerms = append(searchTerms, tools...)

			currentRow = &Row{
				Language:    language,
				Lines:       []string{line},
				SearchTerms: searchTerms,
			}
		} else if len(line) > 0 && currentRow != nil {
			// Continuation line (starts with space/tab)
			currentRow.Lines = append(currentRow.Lines, line)

			// Extract tools from continuation line
			tools := extractTools(line)
			currentRow.SearchTerms = append(currentRow.SearchTerms, tools...)
		}

		i++
	}

	if currentRow != nil {
		rows = append(rows, *currentRow)
	}

	return rows, header, configLines
}

// extractTools extracts tool names from a line (after ✓ or ✘ symbols)
func extractTools(line string) []string {
	var tools []string
	parts := strings.Fields(line)

	for i, part := range parts {
		// Look for ✓ or ✘ followed by a tool name
		if (part == "✓" || part == "✘") && i+1 < len(parts) {
			tool := parts[i+1]
			if tool != "None" {
				// Remove trailing … if present
				tool = strings.TrimSuffix(tool, "…")
				tools = append(tools, tool)
			}
		}
	}

	return tools
}
