package tui

import (
	"assessment-tool-cli/parser"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/selection"
)

type sessionState int

const (
	gradesView sessionState = iota
	topicsView
	questionsView
	editQuestionView
)

func InitGradesModel(data *parser.TOMLData) *GradesModel {
	return &GradesModel{state: gradesView, grades: data.Grades, help: help.New()}
}

type GradesModel struct {
	state     sessionState
	topic     tea.Model
	grades    []parser.GradeData
	selection *selection.Model[string]
	help      help.Model
}

func (m *GradesModel) Init() tea.Cmd {
	grades := make([]string, len(m.grades))
	for i, data := range m.grades {
		grades[i] = data.Grade
	}
	sel := selection.New("Choose grade", grades)
	m.selection = selection.NewModel(sel)
	m.selection.PageSize = len(grades)
	m.selection.Filter = nil

	return tea.Batch(m.selection.Init(), tea.EnterAltScreen)
}

func (m *GradesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.state = topicsView
			c, err := m.selection.ValueAsChoice()
			if err != nil {
				log.Fatal(err)
				return m, tea.Quit
			}

			m.topic = InitTopicsModel(c.String, m.grades[c.Index()].Topics)
			return m.topic, m.topic.Init()
		case "ctrl+c":
			return m, tea.Quit
		default:
			_, cmd = m.selection.Update(msg)
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *GradesModel) View() string {
	switch m.state {
	case topicsView:
		return m.topic.View()
	default:
		var str strings.Builder
		str.WriteString(windowStyle.Width(lipgloss.Width(m.selection.View()) - windowStyle.GetHorizontalFrameSize()).
			Render(m.selection.View() + "\n\n" + helpViewStyle.Render("[↑/↓] - select ◉ [enter] - choose\n[ctrl+c] - quit")))
		return docStyle.Render(str.String())
	}
}
