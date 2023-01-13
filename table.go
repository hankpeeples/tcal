package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
)

var (
	primaryFg = "#555"

	baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(primaryFg)).Align(lipgloss.Left)
)

type model struct {
	calTable table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.calTable.Focused() {
				m.calTable.Blur()
			} else {
				m.calTable.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a":
			return m, tea.Batch(
				tea.Println(addNewCalendarItem()),
			)
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.calTable.SelectedRow()[0]),
			)
		}
	}
	m.calTable, cmd = m.calTable.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.calTable.View()) + "\n"
}

func newModel() model {
	return model{}
}

func printEventList(list []calEvent) tea.Model {
	t := populateTable(list)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(primaryFg)).BorderBottom(true).
		Bold(true)

	t.SetStyles(s)

	tProg, err := tea.NewProgram(model{t}).Run()
	if err != nil {
		pterm.Error.Println("There was an error displaying the calendar table.")
		Log.Fatalf("error running table display: %v", err)
	}

	return tProg
}

func addNewCalendarItem() string {
	// fmt.Println("Add new cal item")
	return "Test return"
}

func populateTable(list []calEvent) table.Model {
	width := pterm.GetTerminalWidth()
	var columns []table.Column

	if width < 100 {
		columns = []table.Column{
			{Title: "Event Name", Width: width / 4},
			{Title: "Date", Width: width / 3},
			{Title: "Last Updated", Width: width / 4},
		}
	} else {
		columns = []table.Column{
			{Title: "Event Name", Width: 25},
			{Title: "Date", Width: 35},
			{Title: "Last Updated", Width: 20},
		}
	}

	var rows []table.Row
	for _, e := range list {
		rows = append(rows, table.Row{e.Name, e.Date, e.Updated})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(list)),
	)

	return t
}
