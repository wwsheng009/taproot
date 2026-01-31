package list

import (
	"strings"
)

// Filter maintains filter state and performs filtering operations.
type Filter struct {
	query    string
	matches  []int
	active   bool
	caseSensitive bool
}

// NewFilter creates a new filter.
func NewFilter() *Filter {
	return &Filter{
		query:    "",
		matches:  nil,
		active:   false,
		caseSensitive: false,
	}
}

// SetQuery sets the filter query.
func (f *Filter) SetQuery(query string) {
	f.query = query
	f.active = query != ""
}

// Query returns the current filter query.
func (f *Filter) Query() string {
	return f.query
}

// Active returns true if filtering is active.
func (f *Filter) Active() bool {
	return f.active
}

// SetCaseSensitive sets whether filtering is case-sensitive.
func (f *Filter) SetCaseSensitive(caseSensitive bool) {
	f.caseSensitive = caseSensitive
}

// Clear clears the filter state.
func (f *Filter) Clear() {
	f.query = ""
	f.matches = nil
	f.active = false
}

// Matches returns the match indexes.
func (f *Filter) Matches() []int {
	return f.matches
}

// SetMatches sets the match indexes.
func (f *Filter) SetMatches(matches []int) {
	f.matches = matches
}

// Apply applies the filter to a list of filterable items.
func (f *Filter) Apply(items []FilterableItem) []FilterableItem {
	if !f.active {
		f.matches = nil
		return items
	}

	var result []FilterableItem
	query := f.prepareQuery()

	for i, item := range items {
		value := item.FilterValue()
		searchValue := value
		if !f.caseSensitive {
			searchValue = strings.ToLower(value)
		}

		if idx := strings.Index(searchValue, query); idx != -1 {
			result = append(result, item)
			// Track which original items matched
			if f.matches == nil {
				f.matches = []int{i}
			} else {
				f.matches = append(f.matches, i)
			}
			// Find all match positions for this item and set them
			itemMatches := f.findAllMatches(value)
			if hm, ok := item.(HasMatchIndexes); ok {
				hm.MatchIndexes(itemMatches)
			}
		}
	}

	return result
}

// ApplyToStrings applies the filter to a list of strings.
func (f *Filter) ApplyToStrings(items []string) []string {
	if !f.active {
		return items
	}

	var result []string
	query := f.prepareQuery()

	for _, item := range items {
		searchValue := item
		if !f.caseSensitive {
			searchValue = strings.ToLower(item)
		}

		if strings.Contains(searchValue, query) {
			result = append(result, item)
		}
	}

	return result
}

// MatchIndexes returns the indexes of items that match the filter.
func (f *Filter) MatchIndexes(items []FilterableItem) []int {
	if !f.active {
		return nil
	}

	var indexes []int
	query := f.prepareQuery()

	for i, item := range items {
		value := item.FilterValue()
		searchValue := value
		if !f.caseSensitive {
			searchValue = strings.ToLower(value)
		}

		if strings.Contains(searchValue, query) {
			indexes = append(indexes, i)
		}
	}

	return indexes
}

// prepareQuery prepares the query string for matching.
func (f *Filter) prepareQuery() string {
	query := f.query
	if !f.caseSensitive {
		query = strings.ToLower(query)
	}
	return query
}

// findAllMatches finds all match positions in a string.
func (f *Filter) findAllMatches(text string) []int {
	var matches []int
	query := f.prepareQuery()
	searchText := text
	if !f.caseSensitive {
		searchText = strings.ToLower(text)
	}

	for {
		idx := strings.Index(searchText, query)
		if idx == -1 {
			break
		}
		matches = append(matches, idx)
		searchText = searchText[idx+len(query):]
	}

	return matches
}

// Highlight returns the text with highlighted matches using ANSI codes.
func (f *Filter) Highlight(text string, before, after string) string {
	if !f.active {
		return text
	}

	query := f.prepareQuery()
	searchText := text
	if !f.caseSensitive {
		searchText = strings.ToLower(text)
	}

	var result strings.Builder
	lastIdx := 0

	for {
		idx := strings.Index(searchText, query)
		if idx == -1 {
			result.WriteString(text[lastIdx:])
			break
		}

		// Add text before match
		result.WriteString(text[lastIdx : lastIdx+idx])
		// Add highlighted match
		result.WriteString(before)
		result.WriteString(text[lastIdx+idx : lastIdx+idx+len(f.query)])
		result.WriteString(after)

		// Move past this match
		lastIdx += idx + len(f.query)
		searchText = searchText[idx+len(query):]
	}

	return result.String()
}

// MatchCount returns the number of matches in a string.
func (f *Filter) MatchCount(text string) int {
	if !f.active {
		return 0
	}

	query := f.prepareQuery()
	searchText := text
	if !f.caseSensitive {
		searchText = strings.ToLower(text)
	}

	count := 0
	for {
		idx := strings.Index(searchText, query)
		if idx == -1 {
			break
		}
		count++
		searchText = searchText[idx+len(query):]
	}

	return count
}

// HasMatchIn returns true if the filter matches any part of the text.
func (f *Filter) HasMatchIn(text string) bool {
	if !f.active {
		return true // No filter means everything matches
	}

	query := f.prepareQuery()
	searchText := text
	if !f.caseSensitive {
		searchText = strings.ToLower(text)
	}

	return strings.Contains(searchText, query)
}
