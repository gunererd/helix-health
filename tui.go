package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Model holds the TUI state
type Model struct {
	textInput    textinput.Model
	viewport     viewport.Model
	rows         []Row
	filteredRows []Row
	header       string   // Original header from helix --health
	configLines  []string // Config info lines from helix --health
	searchQuery  string   // Current search query for highlighting
	ready        bool
	isSearching  bool
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		verticalMarginHeight := tuiHeaderHeight + tuiFooterHeight

		if !m.ready {
			// Initialize viewport on first window size message
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = tuiHeaderHeight
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		// Update viewport content
		m.viewport.SetContent(m.renderRows())

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		case "ctrl+p":
			if m.ready {
				m.viewport.ScrollUp(1)
			}
		case "ctrl+n":
			if m.ready {
				m.viewport.ScrollDown(1)
			}
		case "up", "down", "pgup", "pgdown":
			if m.ready {
				m.viewport, cmd = m.viewport.Update(msg)
				cmds = append(cmds, cmd)
			}
		default:
			// Update text input for other keys
			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)

			// Filter rows based on search input
			searchTerm := m.textInput.Value()
			m.searchQuery = searchTerm
			if searchTerm == "" {
				m.filteredRows = m.rows
				m.isSearching = false
			} else {
				m.filteredRows = make([]Row, 0)
				for _, row := range m.rows {
					if row.matchesSearchTerm(searchTerm) {
						m.filteredRows = append(m.filteredRows, row)
					}
				}
				m.isSearching = true
			}

			// Update viewport content after filtering
			if m.ready {
				m.viewport.SetContent(m.renderRows())
			}
		}
	}

	// Update viewport for other messages
	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// renderRows returns the rows content for the viewport
func (m Model) renderRows() string {
	var s strings.Builder

	for i, row := range m.filteredRows {
		// Add separator between rows only when searching
		if m.isSearching && i > 0 {
			separator := strings.Repeat("─", m.viewport.Width)
			s.WriteString(separatorStyle.Render(separator))
			s.WriteString("\n")
		}
		for _, line := range row.Lines {
			s.WriteString(" ")
			processedLine := line
			// Apply highlighting if searching
			if m.isSearching && m.searchQuery != "" {
				processedLine = highlightLineMatches(line, []string{m.searchQuery})
			}
			processedLine = colorizeSymbols(processedLine)
			s.WriteString(processedLine)
			s.WriteString("\n")
		}
	}

	return s.String()
}

// View renders the TUI
func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	var s strings.Builder

	// Search input
	s.WriteString(" ")
	s.WriteString(searchLabelStyle.Render("Search: "))
	s.WriteString(m.textInput.View())
	s.WriteString("\n\n")

	// Config lines
	for _, line := range m.configLines {
		s.WriteString(" ")
		s.WriteString(line)
		s.WriteString("\n")
	}
	s.WriteString("\n")

	// Column header
	s.WriteString(" ")
	s.WriteString(headerStyle.Render(m.header))
	s.WriteString("\n\n")

	// Rows
	s.WriteString(m.viewport.View())

	// Footer
	scrollPercent := int(m.viewport.ScrollPercent() * 100)
	footer := fmt.Sprintf("%d matches | %d%% | ↑↓ scroll | q/Esc to quit", len(m.filteredRows), scrollPercent)
	s.WriteString("\n ")
	s.WriteString(footerStyle.Render(footer))

	return s.String()
}
