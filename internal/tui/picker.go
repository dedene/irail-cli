package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dedene/irail-cli/internal/api"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#60a5fa"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22c55e")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d4d4d4"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#737373"))

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60a5fa")).
			Bold(true)
)

// PickerResult contains the result of a station picker.
type PickerResult struct {
	Station  api.Station
	Canceled bool
}

type model struct {
	stations     []api.Station
	filtered     []api.Station
	cursor       int
	search       string
	selected     *api.Station
	canceled     bool
	windowHeight int
}

func initialModel(stations []api.Station) model {
	return model{
		stations:     stations,
		filtered:     stations,
		cursor:       0,
		windowHeight: 10,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.canceled = true

			return m, tea.Quit

		case "enter":
			if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
				m.selected = &m.filtered[m.cursor]
			}

			return m, tea.Quit

		case "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "ctrl+n":
			if len(m.filtered) > 0 && m.cursor < len(m.filtered)-1 {
				m.cursor++
			}

		case "backspace":
			if len(m.search) > 0 {
				m.search = m.search[:len(m.search)-1]
				m.filterStations()
			}

		default:
			if len(msg.String()) == 1 {
				m.search += msg.String()
				m.filterStations()
			}
		}

	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height - 4
		if m.windowHeight < 5 {
			m.windowHeight = 5
		}
	}

	return m, nil
}

func (m *model) filterStations() {
	if m.search == "" {
		m.filtered = m.stations
		m.cursor = 0

		return
	}

	search := strings.ToLower(m.search)
	var filtered []api.Station

	for _, s := range m.stations {
		name := strings.ToLower(s.Name)
		stdName := strings.ToLower(s.StandardName)

		if strings.Contains(name, search) || strings.Contains(stdName, search) {
			filtered = append(filtered, s)
		}
	}

	m.filtered = filtered
	m.cursor = 0
}

func (m model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("Select station"))
	b.WriteString("\n")

	// Search input
	searchDisplay := m.search
	if searchDisplay == "" {
		searchDisplay = dimStyle.Render("type to search...")
	} else {
		searchDisplay = inputStyle.Render(searchDisplay)
	}

	b.WriteString(fmt.Sprintf("> %s", searchDisplay))
	b.WriteString("\n\n")

	// Station list
	if len(m.filtered) == 0 {
		b.WriteString(dimStyle.Render("  No matching stations"))
		b.WriteString("\n")
	} else {
		// Calculate visible range
		start := 0
		end := len(m.filtered)

		if end > m.windowHeight {
			// Center the cursor in the visible range
			halfWindow := m.windowHeight / 2

			start = m.cursor - halfWindow
			if start < 0 {
				start = 0
			}

			end = start + m.windowHeight
			if end > len(m.filtered) {
				end = len(m.filtered)
				start = end - m.windowHeight

				if start < 0 {
					start = 0
				}
			}
		}

		for i := start; i < end; i++ {
			s := m.filtered[i]

			cursor := "  "
			style := normalStyle

			if i == m.cursor {
				cursor = "> "
				style = selectedStyle
			}

			line := fmt.Sprintf("%s%s", cursor, style.Render(s.Name))

			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("↑/↓: navigate • enter: select • esc: cancel"))

	return b.String()
}

// PickStation shows an interactive station picker.
func PickStation(stations []api.Station) (PickerResult, error) {
	m := initialModel(stations)

	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return PickerResult{}, fmt.Errorf("run picker: %w", err)
	}

	fm := finalModel.(model)

	if fm.canceled {
		return PickerResult{Canceled: true}, nil
	}

	if fm.selected != nil {
		return PickerResult{Station: *fm.selected}, nil
	}

	return PickerResult{Canceled: true}, nil
}
