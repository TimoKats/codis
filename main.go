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

func QueryStyle(color int) *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color(strconv.Itoa(color)) // 10, 11, 12
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(90)
	return s
}


func ResultStyle() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("0") 
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(90).Height(15)
	return s
}

type model struct {
	query Query
	queryIndex int
	resultIndex int
	width int
	height int
	queryField textinput.Model
	resultField textarea.Model
	queryStyle *Styles
	resultStyle *Styles
}

type Query struct {
	query string
	result []string
	resultLocations []string
	queryType []string
}

func New(query Query) *model {
	queryStyle := QueryStyle(11)
	resultStyle := ResultStyle()
	queryField := textinput.New() 
	queryField.Placeholder = "press / to start querying..."
	resultField := textarea.New()
	resultField.SetWidth(80)
	resultField.SetHeight(15)
	resultField.ShowLineNumbers = false
	return &model{queryIndex: 0, resultIndex:0, query: query, 
		queryField: queryField, resultField: resultField, 
		queryStyle: queryStyle, resultStyle: resultStyle,
	}
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
					m.queryField.Focus()
					return m, nil
				case "esc":
					m.queryField.Reset()
					m.resultField.Reset()
					return m, nil
				case "enter":
					m.query.query = m.queryField.Value()
					m.queryField.Reset()
					m.query.result, m.query.resultLocations = cosearch.BasicQuery(m.query.query, labeledRows, orderedKeys)
					m.resultField.SetValue(m.query.result[m.resultIndex])
					return m, nil
				case "ctrl+j":
					m.queryField.Reset()
					m.queryIndex = (m.queryIndex + 1) % 2
					m.queryStyle = QueryStyle(((m.queryIndex + 1) % 2) + 10)
				case "ctrl+k":
					m.queryField.Reset()
					m.resultField.Reset()
					m.resultIndex = (m.resultIndex + 1) % len(m.query.result)
					m.resultField.SetValue(m.query.result[m.resultIndex])
			}
	}
	m.queryField, cmd = m.queryField.Update(msg)
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
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.query.queryType[m.queryIndex],
				m.queryStyle.InputField.Render(m.queryField.View()),
			),
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.resultStyle.InputField.Render(m.resultField.View()),
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					m.query.resultLocations[m.resultIndex],
					" | ",
					strconv.Itoa(m.resultIndex+1),
					"/",
					strconv.Itoa(len(m.query.result)),
					" | (press ctrl+c to quit)",
				),
			),
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
	query := Query{"", []string{"None"}, []string{"None"}, []string{"basic search", "explorative seach"}}
	m := New(query)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

