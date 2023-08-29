package tui

import (
	"assessment-tool-cli/parser"
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/selection"
)

func InitTopicsModel(grade string, topics []parser.TopicData) *TopicsModel {
	return &TopicsModel{state: topicsView, grade: grade, topics: topics}
}

type TopicsModel struct {
	state     sessionState
	question  tea.Model
	grade     string
	topics    []parser.TopicData
	selection *selection.Model[string]
}

func (m *TopicsModel) Init() tea.Cmd {
	topics := make([]string, len(m.topics))
	for i, data := range m.topics {
		topics[i] = data.Title
	}
	sel := selection.New("Choose topic", topics)
	m.selection = selection.NewModel(sel)
	m.selection.PageSize = len(topics)
	m.selection.Filter = nil

	return m.selection.Init()
}

func (m *TopicsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case QuestionsBackMsg:
		m.state = topicsView
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.state = questionsView
			c, err := m.selection.ValueAsChoice()
			if err != nil {
				log.Fatal(err)
				return m, tea.Quit
			}

			m.question = InitQuestionsModel(m, c.String, m.topics[c.Index()].Questions)
			return m.question, m.question.Init()
		case "ctrl+c":
			parser.EncodeTOML(parser.GradeData{Grade: m.grade, Topics: m.topics})
			return m, tea.Quit
		default:
			_, cmd = m.selection.Update(msg)
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *TopicsModel) View() string {
	switch m.state {
	case questionsView:
		return m.question.View()
	default:
		var str strings.Builder
		str.WriteString(windowStyle.Width(lipgloss.Width(m.selection.View()) - windowStyle.GetHorizontalFrameSize()).
			Render(fmt.Sprintf("Grade - %v\n\n", m.grade) + m.selection.View() + "\n\n" +
				helpViewStyle.Render("[↑/↓] - select ◉ [enter] - choose\n[ctrl+c] - quit and save all data")))
		return docStyle.Render(str.String())
	}
}
