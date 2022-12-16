package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#c0a46b"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func printEventList(list []calEvent) {
	columns := []table.Column{
		{Title: "Event Name", Width: pterm.GetTerminalWidth() / 7},
		{Title: "Date", Width: pterm.GetTerminalWidth() / 7},
		{Title: "Description", Width: pterm.GetTerminalWidth() / 7},
		{Title: "Type", Width: pterm.GetTerminalWidth() / 7},
		{Title: "Status", Width: pterm.GetTerminalWidth() / 7},
		{Title: "Updated", Width: pterm.GetTerminalWidth() / 7},
	}

	var rows []table.Row
	for _, e := range list {
		rows = append(rows, table.Row{e.Name, e.Date, e.Description, e.Type, e.Status, e.Updated})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(list)),
	)

	tea.NewProgram(model{t}).Run()
}
