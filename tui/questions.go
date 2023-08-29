package tui

import (
	"assessment-tool-cli/parser"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/muesli/termenv"
)

type QuestionsBackMsg bool

var (
	topicToQuestions = make(map[string][]string)
	tabs             = []string{"Question", "Answer", "Feedback"}
)

func InitQuestionsModel(topicModel tea.Model, topic string, fullQuestions []parser.QuestionData) *QuestionsModel {
	_, ok := topicToQuestions[topic]
	if !ok {
		questions := make([]string, len(fullQuestions))
		for i, data := range fullQuestions {
			questions[i] = data.Question
		}
		topicToQuestions[topic] = questions
	}

	return &QuestionsModel{state: questionsView, topicModel: topicModel, topic: topic, fullQuestions: fullQuestions}
}

type QuestionsModel struct {
	state         sessionState
	topic         string
	topicModel    tea.Model
	fullQuestions []parser.QuestionData
	selection     *selection.Model[string]
	areas         []textarea.Model
	questionIndex int
	activeArea    int
}

func (m *QuestionsModel) Init() tea.Cmd {
	sel := selection.New("Choose question", topicToQuestions[m.topic])
	m.selection = selection.NewModel(sel)
	m.selection.PageSize = len(topicToQuestions[m.topic])
	m.selection.Filter = nil

	m.areas = make([]textarea.Model, 3)
	for i := range m.areas {
		m.areas[i] = textarea.New()
	}

	return m.selection.Init()
}

func (m *QuestionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			var previousArea int
			if m.activeArea == len(m.areas)-1 {
				m.activeArea = 0
				previousArea = 2
			} else {
				m.activeArea++
				previousArea = m.activeArea - 1
			}
			m.areas[m.activeArea].Focus()
			m.areas[m.activeArea].FocusedStyle.Text = focusedStyle
			m.areas[m.activeArea].FocusedStyle.Prompt = focusedStyle
			m.areas[previousArea].Blur()
			m.areas[previousArea].FocusedStyle.Text = noStyle
			m.areas[previousArea].FocusedStyle.Prompt = noStyle
		case "enter":
			m.state = editQuestionView
			c, err := m.selection.ValueAsChoice()
			if err != nil {
				log.Fatal(err)
				return m, tea.Quit
			}
			m.questionIndex = c.Index()
			m.areas[0].InsertString(m.fullQuestions[m.questionIndex].Question)
			m.areas[0].Focus()
			m.areas[0].FocusedStyle.Text = focusedStyle
			m.areas[0].FocusedStyle.Prompt = focusedStyle
			m.areas[1].InsertString(m.fullQuestions[m.questionIndex].Answer)
			m.areas[2].InsertString(m.fullQuestions[m.questionIndex].Feedback)
		case "+":
			err := HighLightQuestion(m, true)
			if err != nil {
				log.Fatal(err)
				return m, tea.Quit
			}
			cmd = m.selection.Init()
		case "_":
			err := HighLightQuestion(m, false)
			if err != nil {
				log.Fatal(err)
				return m, tea.Quit
			}
			cmd = m.selection.Init()
		case "esc":
			if m.state == editQuestionView {
				m.state = questionsView
				m.fullQuestions[m.questionIndex].Question = m.areas[0].Value()
				m.fullQuestions[m.questionIndex].Answer = m.areas[1].Value()
				m.fullQuestions[m.questionIndex].Feedback = m.areas[2].Value()
				m.areas = make([]textarea.Model, 3)
				for i := range m.areas {
					m.areas[i] = textarea.New()
				}
				return m, cmd
			}
			return m.topicModel, func() tea.Msg {
				return QuestionsBackMsg(true)
			}
		case "ctrl+c":
			return m.topicModel.Update(msg)
		default:
			if m.state == questionsView {
				_, cmd = m.selection.Update(msg)
			} else {
				cmd = m.updateAreas(msg)
			}
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *QuestionsModel) updateAreas(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.areas))

	for i := range m.areas {
		m.areas[i], cmds[i] = m.areas[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func (m *QuestionsModel) View() string {
	var str strings.Builder
	switch m.state {
	case editQuestionView:
		var renderedTabs []string

		for i, t := range tabs {
			var style lipgloss.Style
			isFirst, isLast, isActive := i == 0, i == len(tabs)-1, i == m.activeArea
			if isActive {
				style = activeTabStyle.Copy()
			} else {
				style = inactiveTabStyle.Copy()
			}
			border, _, _, _, _ := style.GetBorder()
			if isFirst && isActive {
				border.BottomLeft = "│"
			} else if isFirst && !isActive {
				border.BottomLeft = "├"
			} else if isLast && isActive {
				border.BottomRight = "│"
			} else if isLast && !isActive {
				border.BottomRight = "┤"
			}
			style = style.Border(border)
			renderedTabs = append(renderedTabs, style.Render(t))
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
		str.WriteString(row + "\n" + windowStyle.Padding(0).UnsetBorderTop().
			Width(lipgloss.Width(row)-2).Render(m.areas[m.activeArea].View()) + "\n\n" +
			helpViewStyle.Render("\t\t[tab] - switch between ◉ [esc] - return\n\t\t[ctrl+c] - quit and save all data"))
		return docStyle.Render(str.String())
	default:
		str.WriteString(windowStyle.BorderTop(true).Padding(2, 16).
			Width(lipgloss.Width(m.selection.View()) - windowStyle.GetHorizontalFrameSize()).
			Render(fmt.Sprintf("Topic - %v\n\n", m.topic) + m.selection.View() + "\n\n" +
				helpViewStyle.Render("[↑/↓] - select ◉ [enter] - choose ◉ [esc]: return\n"+
					"[+/_] - mark correctly/incorrectly answered\n[ctrl+c] - quit and save all data")))
		return docStyle.Render(str.String())
	}
}

func HighLightQuestion(m *QuestionsModel, done bool) error {
	c, err := m.selection.ValueAsChoice()
	if err != nil {
		return err
	}

	i := c.Index()
	if !done {
		topicToQuestions[m.topic][i] = termenv.String(topicToQuestions[m.topic][i]).Foreground(clr.Color("0")).
			Background(clr.Color("#E88388")).String()
	} else {
		topicToQuestions[m.topic][i] = termenv.String(topicToQuestions[m.topic][i]).Foreground(clr.Color("0")).
			Background(clr.Color("#A8CC8C")).String()
	}
	m.selection.Selection = selection.New("Choose question", topicToQuestions[m.topic])
	m.selection.PageSize = len(topicToQuestions[m.topic])
	m.selection.Filter = nil

	return nil
}
