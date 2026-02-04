package output

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

// ansiRegex matches ANSI color escape sequences.
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// oscRegex matches OSC 8 hyperlink sequences: \x1b]8;;URL\x1b\\ or \x1b]8;;URL\x07
var oscRegex = regexp.MustCompile(`\x1b\]8;;[^\x07\x1b]*(?:\x07|\x1b\\)`)

// Table provides ANSI-aware aligned table output.
type Table struct {
	w         io.Writer
	headers   []string
	rows      [][]string
	colWidths []int
	padding   int
}

// NewTable creates a new table writer.
func NewTable(w io.Writer) *Table {
	return &Table{
		w:       w,
		padding: 2,
	}
}

// SetHeaders sets the table headers.
func (t *Table) SetHeaders(headers ...string) {
	t.headers = headers
	t.updateWidths(headers)
}

// WriteHeader is a no-op; headers are written in Render.
func (t *Table) WriteHeader() {
	// Headers are written during Render
}

// AddRow adds a row to the table.
func (t *Table) AddRow(cols ...string) {
	t.rows = append(t.rows, cols)
	t.updateWidths(cols)
}

// updateWidths updates column widths based on visible text width.
func (t *Table) updateWidths(cols []string) {
	for i, col := range cols {
		w := visibleWidth(col)
		if i >= len(t.colWidths) {
			t.colWidths = append(t.colWidths, w)
		} else if w > t.colWidths[i] {
			t.colWidths[i] = w
		}
	}
}

// visibleWidth returns the visible width of a string, ignoring ANSI codes and hyperlinks.
func visibleWidth(s string) int {
	// Strip OSC 8 hyperlink sequences first
	s = oscRegex.ReplaceAllString(s, "")
	// Strip ANSI color codes
	s = ansiRegex.ReplaceAllString(s, "")
	return len(s)
}

// Render outputs the table.
func (t *Table) Render() error {
	// Write headers
	if len(t.headers) > 0 {
		if err := t.writeRow(t.headers); err != nil {
			return err
		}
	}

	// Write data rows
	for _, row := range t.rows {
		if err := t.writeRow(row); err != nil {
			return err
		}
	}

	return nil
}

// writeRow writes a single row with proper padding.
func (t *Table) writeRow(cols []string) error {
	var sb strings.Builder

	for i, col := range cols {
		sb.WriteString(col)

		// Add padding to reach column width + spacing (except last column)
		if i < len(cols)-1 {
			visible := visibleWidth(col)
			width := 0
			if i < len(t.colWidths) {
				width = t.colWidths[i]
			}

			padding := width - visible + t.padding
			if padding > 0 {
				sb.WriteString(strings.Repeat(" ", padding))
			}
		}
	}

	sb.WriteString("\n")
	_, err := fmt.Fprint(t.w, sb.String())

	return err
}
