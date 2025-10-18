// helix-health provides an interactive TUI for viewing and searching Helix editor's health information.
//
// It wraps the `helix --health` command and offers:
//   - Interactive search with real-time filtering
//   - Non-interactive mode for command-line queries
//   - Syntax highlighting and colored status indicators
//
// Usage:
//
//	helix-health          # Launch interactive TUI
//	helix-health python   # Filter for python
//	helix-health go rust  # Filter for multiple languages
package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// findHelixCommand searches for the helix editor command in PATH.
func findHelixCommand() (string, error) {
	if _, err := exec.LookPath("helix"); err == nil {
		return "helix", nil
	}

	if _, err := exec.LookPath("hx"); err == nil {
		return "hx", nil
	}

	return "", fmt.Errorf("helix editor not found in PATH (tried 'helix' and 'hx')")
}

// Row represents a language row from helix --health
type Row struct {
	Language    string
	Lines       []string
	SearchTerms []string // All searchable terms: [language, tool1, tool2, ...]
}

// matchesSearchTerm checks if any search term in row contains the query (case-insensitive substring)
func (r Row) matchesSearchTerm(query string) bool {
	queryLower := strings.ToLower(query)
	for _, term := range r.SearchTerms {
		if strings.Contains(strings.ToLower(term), queryLower) {
			return true
		}
	}
	return false
}

// UI layout constants
const (
	tuiFooterHeight = 1
)

// Styles
var (
	searchLabelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	headerStyle      = lipgloss.NewStyle().Bold(true).Underline(false).Foreground(lipgloss.Color("14"))
	footerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	separatorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("235"))
	checkStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))             // Green
	xStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))             // Red
	boldStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11")) // Bold yellow
)

func main() {
	helixCmd, err := findHelixCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Please install Helix editor: https://helix-editor.com/")
		os.Exit(1)
	}

	cmd := exec.Command(helixCmd, "--health")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running %s --health: %v\n", helixCmd, err)
		os.Exit(1)
	}

	rows, header, configLines := parseHealthOutput(string(output))

	// Check if arguments are provided (non-interactive mode)
	if len(os.Args) > 1 {
		searchTerms := os.Args[1:]
		runNonInteractive(rows, header, configLines, searchTerms)
		return
	}

	// Interactive TUI mode
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	m := Model{
		textInput:    ti,
		rows:         rows,
		filteredRows: rows,
		header:       header,
		configLines:  configLines,
		ready:        false,
		isSearching:  false,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func runNonInteractive(rows []Row, header string, configLines []string, searchTerms []string) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80 // fallback to 80 if we can't detect
	}

	// Print config lines as-is
	for _, line := range configLines {
		fmt.Println(line)
	}
	fmt.Println()

	// Filter rows
	var filteredRows []Row
	for _, row := range rows {
		if slices.ContainsFunc(searchTerms, row.matchesSearchTerm) {
			filteredRows = append(filteredRows, row)
		}
	}

	// Print header
	fmt.Println(headerStyle.Render(header))

	if len(filteredRows) == 0 {
		fmt.Println("No matches found")
		return
	}

	// Print filtered language rows
	for i, row := range filteredRows {
		// Add separator between rows
		if i > 0 {
			separator := strings.Repeat("─", width)
			fmt.Println(separatorStyle.Render(separator))
		}

		for _, line := range row.Lines {
			processedLine := highlightLineMatches(line, searchTerms)
			processedLine = colorizeSymbols(processedLine)
			fmt.Println(processedLine)
		}
	}
}
