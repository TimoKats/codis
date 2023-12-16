/* 
** @name: main 
** @author: Timo Kats
** @description: Contains the TUI and calls search/parsing functions. 
*/

package main

import (
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	coinit "codis/coinit"
	cosearch "codis/cosearch"
	coexplore "codis/coexplore"
)

// globals

var fullTree, _ = coexplore.NewTree(coinit.CurrentDirectory)

// structs 

type Styles struct {
	BorderColor lipgloss.Color
	InputField lipgloss.Style
}

func QueryStyle(color int) *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color(strconv.Itoa(color)) // 10, 11, 12
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(100)
	return s
}

func ResultStyle() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("0") 
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(100).Height(20)
	return s
}

type model struct {
	query Query
	queryIndex int
	resultIndex int
	infoIndex int
	width int
	height int
	viewDirOnly bool
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

/* 
** @name: New 
** @description: Initiates a new TUI with default values.
*/
func New(query Query) *model {
	queryStyle := QueryStyle(10)
	resultStyle := ResultStyle()
	queryField := textinput.New() 
	queryField.Placeholder = "press / to start querying..."
	resultField := textarea.New()
	resultField.SetWidth(100)
	resultField.SetHeight(20)
	resultField.ShowLineNumbers = false
	resultField.CharLimit = -1
	return &model{queryIndex: 0, resultIndex:0, infoIndex: 0,
	query: query, viewDirOnly: false, queryField: queryField, 
	resultField: resultField, queryStyle: queryStyle, resultStyle: resultStyle,
	}
} 

/* 
** @name: Init
** @description: Mandetory(?) function that starts the TUI.
*/
func (m model) Init() tea.Cmd {
	return nil
}

// keypresses

/* 
** @name: KeyEscape 
** @description: Empties the current set of results. 
*/
func KeyEscape(m model) (tea.Model, tea.Cmd) {
	m.queryField.Reset()
	m.resultField.Reset()
	m.resultIndex = 0
	return m, nil
}

/* 
** @name: KeyEnter 
** @description: Runs the selected query. 
*/
func KeyEnter(m model) (tea.Model, tea.Cmd) {
	m.resultIndex = 0
	m.query.query = m.queryField.Value()
	m.queryField.Reset()
	if m.queryIndex == 0 {
		m.query.result, m.query.resultLocations = cosearch.BasicQuery(m.query.query)
	} else if m.queryIndex == 1 {
		m.query.result, m.query.resultLocations = cosearch.FuzzyQuery(m.query.query)
	} else {
		m.query.result, m.query.resultLocations = coexplore.Show(fullTree, 0, 5, m.query.query, m.viewDirOnly, m.infoIndex)
	}
	m.resultField.SetValue(m.query.result[m.resultIndex])
	return m, nil
} 

/* 
** @name: KeyBack 
** @description: Switches to the previous search result. 
*/
func KeyBack(m model) (tea.Model, tea.Cmd) {
	m.queryField.Reset()
	m.resultField.Reset()
	if m.resultIndex > 0 {
		m.resultIndex = (m.resultIndex - 1) % len(m.query.result)
	} else {
		m.resultIndex = len(m.query.result) - 1
	}
	m.resultField.SetValue(m.query.result[m.resultIndex])
	return m, nil
} 

/* 
** @name: KeyForward 
** @description: Switches to the next search result. 
*/
func KeyForward(m model) (tea.Model, tea.Cmd) {
	m.queryField.Reset()
	m.resultField.Reset()
	m.resultIndex = (m.resultIndex + 1) % len(m.query.result)
	m.resultField.SetValue(m.query.result[m.resultIndex])
	return m, nil
} 

/* 
** @name: Tab
** @description: Switches to the next search function.
*/
func KeyTab(m model) (tea.Model, tea.Cmd) {
	m.queryIndex = (m.queryIndex + 1) % len(m.query.queryType)
	m.queryStyle = QueryStyle((m.queryIndex % 3) + 10)
	return m, nil
}

/* 
** @name: Ctrl+Tab
** @description: Switches to the types of info boxes in explore mode 
*/
func KeyCtrlG(m model) (tea.Model, tea.Cmd) {
	if m.queryIndex == 2 {
		m.infoIndex = (m.infoIndex + 1) % 2 // temporary like this
		m.query.result, m.query.resultLocations = coexplore.Show(fullTree, 0, 5, m.query.query, m.viewDirOnly, m.infoIndex)
		m.resultField.SetValue(m.query.result[m.resultIndex])
	}
	return m, nil
}

/* 
** @name: KeyToggleDir
** @description: Switches between file and directory view 
*/
func KeyToggleDir(m model) (tea.Model, tea.Cmd) {
	if m.queryIndex == 2 {
		m.viewDirOnly = !m.viewDirOnly
		m.query.result, m.query.resultLocations = coexplore.Show(fullTree, 0, 5, m.query.query, m.viewDirOnly, m.infoIndex)
		m.resultField.SetValue(m.query.result[m.resultIndex])
	}
	return m, nil
}

/* 
** @name: Update
** @description: Contains the actions for different keypresses in the TUI. 
*/
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
					return KeyEscape(m)
				case "enter":
					return KeyEnter(m)
				case "ctrl+d":
					return KeyToggleDir(m)
				case "ctrl+j":
					return KeyBack(m)
				case "ctrl+k":
					return KeyForward(m)
				case "tab":
					return KeyTab(m)
				case "ctrl+g":
					return KeyCtrlG(m)
			}
	}
	m.queryField, cmd = m.queryField.Update(msg)
	return m, cmd	
}

/* 
** @name: View
** @description: Returns the layout/placement of the visual elements of the TUI.
*/
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

func main() {
	query := Query{"", []string{"None"}, []string{"None"}, []string{"Quick search", "Fuzzy search", "Explorative search"}}
	m := New(query)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

