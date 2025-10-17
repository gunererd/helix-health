package main

import (
	"strings"
)

// colorizeSymbols replaces status symbols with colored versions
func colorizeSymbols(line string) string {
	// Replace ✓ with green version
	line = strings.ReplaceAll(line, "✓", checkStyle.Render("✓"))
	// Replace ✘ with red version
	line = strings.ReplaceAll(line, "✘", xStyle.Render("✘"))
	return line
}

// highlightMatches highlights substring matches in a string (case-insensitive)
func highlightMatches(text string, searchQuery string) string {
	if searchQuery == "" {
		return text
	}

	// Find substring match (case-insensitive)
	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(searchQuery)

	idx := strings.Index(lowerText, lowerQuery)
	if idx == -1 {
		return text
	}

	// Build highlighted string
	before := text[:idx]
	match := text[idx : idx+len(searchQuery)]
	after := text[idx+len(searchQuery):]

	return before + boldStyle.Render(match) + after
}

// highlightLineMatches highlights matches in a line for given search terms
func highlightLineMatches(line string, searchTerms []string) string {
	if len(searchTerms) == 0 {
		return line
	}

	// Extract words with their positions in the original line
	type wordPos struct {
		word  string
		start int
		end   int
	}

	var words []wordPos
	inWord := false
	wordStart := 0

	runes := []rune(line)
	for i, r := range runes {
		if r == ' ' || r == '\t' {
			if inWord {
				words = append(words, wordPos{
					word:  string(runes[wordStart:i]),
					start: wordStart,
					end:   i,
				})
				inWord = false
			}
		} else {
			if !inWord {
				wordStart = i
				inWord = true
			}
		}
	}
	// Handle last word
	if inWord {
		words = append(words, wordPos{
			word:  string(runes[wordStart:]),
			start: wordStart,
			end:   len(runes),
		})
	}

	// Check which words match and build highlighted versions
	highlightedWords := make(map[int]string) // position -> highlighted word
	for _, wp := range words {
		// Skip highlighting "None"
		if wp.word == "None" {
			continue
		}
		for _, searchTerm := range searchTerms {
			if searchTerm == "" {
				continue
			}
			// Check for substring match (case-insensitive)
			if strings.Contains(strings.ToLower(wp.word), strings.ToLower(searchTerm)) {
				highlightedWords[wp.start] = highlightMatches(wp.word, searchTerm)
				break // Only highlight once per word
			}
		}
	}

	// Build result string with highlights
	var result strings.Builder
	i := 0
	for i < len(runes) {
		if highlighted, ok := highlightedWords[i]; ok {
			result.WriteString(highlighted)
			// Skip to end of this word
			for _, wp := range words {
				if wp.start == i {
					i = wp.end
					break
				}
			}
		} else {
			result.WriteRune(runes[i])
			i++
		}
	}

	return result.String()
}
