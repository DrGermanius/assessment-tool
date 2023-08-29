package tui

import (
	"assessment-tool-cli/parser"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	p                 *tea.Program
	clr               = termenv.ColorProfile()
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(10, 0, 0, 30)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 2).Align(lipgloss.Center)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true).Padding(0, 8)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 16).Align(lipgloss.Center).Border(lipgloss.NormalBorder())
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle           = lipgloss.NewStyle()
	helpViewStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func StartTea(data *parser.TOMLData) {
	m := InitGradesModel(data)
	p = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
