package main

import (
	"fmt"
	"os"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	coparse "codis/coparse"
	cosearch "codis/cosearch"
)

// global declarations

var currentDirectory, _ = os.Getwd()
var labeledRows, orderedKeys = coparse.ReturnLabels(currentDirectory)

// view functions

type Styles struct {
	BorderColor lipgloss.Color
	InputField lipgloss.Style
}

func DefaultStyles(color int) *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color(strconv.Itoa(color)) // 10, 11, 12
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type model struct {
	queryIndex int
	resultIndex int
	query Query
	width int
	height int
	answerField textinput.Model
	resultField textarea.Model
	styles *Styles
}

type Query struct {
	query string
	result []string
	queryType []string
}

func New(query Query) *model {
	styles := DefaultStyles(11)
	answerField := textinput.New() 
	answerField.Placeholder = "press / to start querying..."
	resultField := textarea.New()
	return &model{queryIndex: 0, resultIndex:0, query: query, answerField: answerField, resultField: resultField, styles: styles}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
		case tea.KeyMsg:
			switch msg.String() {
				case "ctrl+c":
					return m, tea.Quit
				case "/":
					m.answerField.Focus()
					return m, nil
				case "esc":
					m.answerField.Reset()
					m.resultField.Reset()
					return m, nil
				case "enter":
					m.query.query = m.answerField.Value()
					m.answerField.Reset()
					m.query.result = cosearch.BasicQuery(m.query.query, labeledRows, orderedKeys)
					m.resultField.SetValue(m.query.result[m.resultIndex])
					return m, nil
				case ">":
					m.answerField.Reset()
					m.queryIndex = (m.queryIndex + 1) % 2
					m.styles = DefaultStyles(((m.queryIndex + 1) % 2) + 10)
				case "n":
					m.answerField.Reset()
					m.resultField.Reset()
					m.resultIndex = (m.resultIndex + 1) % len(m.query.result)
					m.resultField.SetValue(m.query.result[m.resultIndex])
			}
	}
	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd	
}

func (m model) View() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Center, 
			m.query.queryType[m.queryIndex], 
			m.styles.InputField.Render(m.answerField.View()),
			m.resultField.View(),
		),
	)
}

// other functions

func printLabels(labeledRows map[coparse.RowLabel]string, orderedKeys []coparse.RowLabel) {
	for _, key := range orderedKeys {
		fmt.Println(key, labeledRows[key])
	}
}

func tempCodis() {
	currentDirectory, _ := os.Getwd()
	labeledRows, orderedKeys := coparse.ReturnLabels(currentDirectory)
	items := cosearch.GetFileTypes(labeledRows, orderedKeys)
	for key, value := range items {
		fmt.Println(key, value)
	}	
	items = cosearch.GetFileCategories(labeledRows, orderedKeys)
	for key, value := range items {
		fmt.Println(key, value)
	}
}

func main() {
	tempCodis()
	query := Query{"", []string{}, []string{"basic search", "explorative seach"}}
	m := New(query)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

