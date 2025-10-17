package main

import (
	"bufio"
	"strings"
)

// parseHealthOutput parses helix --health output into rows.
// Returns rows, column header, and config info lines.
func parseHealthOutput(output string) ([]Row, string, []string) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	var rows []Row
	var currentRow *Row
	var header string
	var configLines []string
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Capture config info lines
		if lineNum <= healthOutputConfigLines {
			configLines = append(configLines, line)
			continue
		}

		// Capture the column header line
		if lineNum == healthOutputHeaderLine {
			header = line
			continue
		}

		// Check if this is a new row (line doesn't start with space or tab)
		if len(line) > 0 && line[0] != ' ' && line[0] != '\t' {
			// Save previous row
			if currentRow != nil {
				rows = append(rows, *currentRow)
			}

			// Start new row
			parts := strings.Fields(line)
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
