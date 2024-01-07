/* 
** @name: main 
** @author: Timo Kats
** @description: Contains the TUI and calls search/parsing functions. 
*/

package main

import (
	"fmt"
	"log"
	"strings"
	"strconv"
	
	// only external dependencies
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"

	coparse "codis/lib/coparse"
	coutils "codis/lib/coutils"
	cotypes "codis/lib/cotypes"
	cofile "codis/lib/cofile"
	cosearch "codis/lib/cosearch"
	coexplore "codis/lib/coexplore"
	cocommands "codis/lib/cocommands"
	codependencies "codis/lib/codependencies"
)

// globals that need to remain constant

var fullTree, _ = coexplore.NewTree(coparse.CurrentDirectory)
var rootFiles = codependencies.GetRootFiles()

// structs 

type Styles struct {
	BorderColor lipgloss.Color
	InputField lipgloss.Style
}

func QueryStyle(color int) *Styles {
	s:= new(Styles)
	s.BorderColor = lipgloss.Color(strconv.Itoa(color))
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
	indecies cotypes.Indecies
	query cotypes.Query
	width int
	height int 
	viewDirOnly bool
	commandMode bool
	formMode bool
	queryField textinput.Model
	resultField textarea.Model
	queryStyle *Styles
	resultStyle *Styles
	contextCategories []int
	contextComment []int
}

/* 
** @name: New 
** @description: Initiates a new TUI with default values.
*/
func New(query cotypes.Query) *model {
	indecies := cotypes.Indecies{QueryIndex:0,ResultIndex:0,InfoIndex:0,FormIndex:1,FileViewIndex:0}
	queryStyle := QueryStyle(10)
	resultStyle := ResultStyle()
	queryField := textinput.New() 
	queryField.Placeholder = "press / to start querying or : for command mode"
	resultField := textarea.New()
	resultField.SetWidth(100)
	resultField.SetHeight(20)
	resultField.ShowLineNumbers = false
	resultField.CharLimit = -1
	return &model{formMode: false, indecies: indecies, query: query, 
	viewDirOnly: false, commandMode: false, queryField: queryField, 
	resultField: resultField, queryStyle: queryStyle, resultStyle: resultStyle,
	contextCategories: []int{}, contextComment: []int{1,0},
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
	m.indecies.ResultIndex = 0
	return m, nil
}

func KeyEnterSearch (m model) (tea.Model, tea.Cmd) {
	categoryContext := coutils.SubsetSlice(coparse.ContextCategories, m.contextCategories)
	infoContext := bool(m.contextComment[0] == 1)
	if m.indecies.QueryIndex == 0 { 
		m.query.Result, m.query.ResultLocations = cosearch.BasicQuery(m.query.Query, categoryContext, infoContext)
	} else if m.indecies.QueryIndex == 1 {
		m.query.Result, m.query.ResultLocations = cosearch.FuzzyQuery(m.query.Query, categoryContext, infoContext)
	} else if m.indecies.QueryIndex == 2 {
		m.query.Result, m.query.ResultLocations = coexplore.Show(fullTree, 0, 5, m.query.Query, m.viewDirOnly, m.indecies.InfoIndex)
	} else if m.indecies.QueryIndex == 3 {
		m.query.Result, m.query.ResultLocations = codependencies.Show(m.indecies.InfoIndex, rootFiles, m.query.Query)
	} else if m.indecies.QueryIndex == 4 {
		m.query.Result, m.query.ResultLocations = cofile.Show(m.query.Query, m.indecies.FileViewIndex, categoryContext)
	}	
	m.resultField.SetValue(m.query.Result[m.indecies.ResultIndex])
	return m, nil
}

func KeyEnterForm(m model) (tea.Model, tea.Cmd) {
	if m.indecies.ContextIndex == 0 {
		if !coutils.ContainsInt(m.contextCategories, m.indecies.FormIndex) {
			m.contextCategories = append(m.contextCategories, m.indecies.FormIndex)
		} else {
			m.contextCategories = coutils.DeleteInt(m.contextCategories, m.indecies.FormIndex)
		}
	} else if m.indecies.ContextIndex == 1 {
		if m.contextComment[0] == 0 {
			m.contextComment = []int{1,0}
		} else {
			m.contextComment = []int{0,1} 
		}
	}
	return m, nil
}

func KeyEnterCommand(m model) (tea.Model, tea.Cmd) {
	m.query.Result, m.query.ResultLocations = cocommands.ParseCommand(m.query.Query)
	m.resultField.SetValue(m.query.Result[m.indecies.ResultIndex])
	return m, nil
}

/* 
** @name: KeyEnter 
** @description: Runs the selected query. 
*/
func KeyEnter(m model) (tea.Model, tea.Cmd) {
	m.indecies.ResultIndex = 0
	m.query.Query = m.queryField.Value()
	m.queryField.Reset()
	if !m.commandMode && !m.formMode {
		return KeyEnterSearch(m)
	} else if m.commandMode && !m.formMode {
		return KeyEnterCommand(m)
	} else if !m.commandMode && m.formMode {
		return KeyEnterForm(m)
	}
	return m, nil
} 

/* 
** @name: KeyBack 
** @description: Switches to the previous search result. 
*/
func KeyBack(m model) (tea.Model, tea.Cmd) {
	m.queryField.Reset()
	m.resultField.Reset()
	if m.indecies.ResultIndex > 0 {
		m.indecies.ResultIndex = (m.indecies.ResultIndex - 1) % len(m.query.Result)
	} else {
		m.indecies.ResultIndex = len(m.query.Result) - 1
	}
	m.resultField.SetValue(m.query.Result[m.indecies.ResultIndex])
	return m, nil
} 

/* 
** @name: KeyForward 
** @description: Switches to the next search result. 
*/
func KeyForward(m model) (tea.Model, tea.Cmd) {
	m.queryField.Reset()
	m.resultField.Reset()
	m.indecies.ResultIndex = (m.indecies.ResultIndex + 1) % len(m.query.Result)
	m.resultField.SetValue(m.query.Result[m.indecies.ResultIndex])
	return m, nil
} 

/* 
** @name: Tab
** @description: Switches to the next search function.
*/
func KeyTab(m model) (tea.Model, tea.Cmd) {
	if !m.commandMode && !m.formMode {
		m.indecies.QueryIndex = (m.indecies.QueryIndex + 1) % len(m.query.QueryType)
		m.queryStyle = QueryStyle((m.indecies.QueryIndex % 5) + 10)
	} else if m.formMode {
		m.indecies.ContextIndex = (m.indecies.ContextIndex + 1) % 2 
	}
	return m, nil
}

/* 
** @name: Ctrl+Tab
** @description: Switches to the types of info boxes in explore mode 
*/
func KeyCtrlG(m model) (tea.Model, tea.Cmd) {
	categoryContext := coutils.SubsetSlice(coparse.ContextCategories, m.contextCategories)
	if m.indecies.QueryIndex == 2 {
		m.indecies.InfoIndex = (m.indecies.InfoIndex + 1) % len(coparse.InfoBoxCategories) 
		m.query.Result, m.query.ResultLocations = coexplore.Show(fullTree, 0, 5, m.query.Query, m.viewDirOnly, m.indecies.InfoIndex)
		m.resultField.SetValue(m.query.Result[m.indecies.ResultIndex])
	} else if m.indecies.QueryIndex == 3 {
		m.indecies.InfoIndex = (m.indecies.InfoIndex + 1) % len(coparse.InfoBoxCategories)
		m.query.Result, m.query.ResultLocations = codependencies.Show(m.indecies.InfoIndex, rootFiles, m.query.Query)
		m.resultField.SetValue(m.query.Result[m.indecies.ResultIndex])
	} else if m.indecies.QueryIndex == 4 {
		m.indecies.FileViewIndex = (m.indecies.FileViewIndex + 1) % 2 
		m.query.Result, m.query.ResultLocations = cofile.Show(m.query.Query, m.indecies.FileViewIndex, categoryContext)
		m.resultField.SetValue(m.query.Result[m.indecies.ResultIndex])
	}
	return m, nil
}

/* 
** @name: KeyToggleDir
** @description: Switches between file and directory view 
*/
func KeyToggleDir(m model) (tea.Model, tea.Cmd) {
	if m.indecies.QueryIndex == 2 {
		m.viewDirOnly = !m.viewDirOnly
		m.query.Result, m.query.ResultLocations = coexplore.Show(fullTree, 0, 5, m.query.Query, m.viewDirOnly, m.indecies.InfoIndex)
		m.resultField.SetValue(m.query.Result[m.indecies.ResultIndex])
	}
	return m, nil
}

func KeyCtrlF(m model) (tea.Model, tea.Cmd) {
	m.formMode = !m.formMode
	return m, nil
}

func KeyUp(m model) (tea.Model, tea.Cmd) {
	if m.formMode { 
		m.indecies.FormIndex = (m.indecies.FormIndex + 1) % len(coparse.ContextCategories)
	}
	return m, nil
}

func KeyDown(m model) (tea.Model, tea.Cmd) {
	if m.formMode { 
		m.indecies.FormIndex = (m.indecies.FormIndex - 1) % len(coparse.ContextCategories)
		if m.indecies.FormIndex < 0 {
			m.indecies.FormIndex = len(coparse.ContextCategories) - 1
		}
	}
	return m, nil
}

/* 
** @name: KeyColon
** @description: Switches to command mode 
*/
func KeyColon(m model) (tea.Model, tea.Cmd) {
	m.indecies.ResultIndex = 0
	m.query.Result, m.query.ResultLocations = []string{""}, []string{"None"} 
	m.resultField.SetValue(m.query.Result[m.indecies.ResultIndex])
	m.commandMode = !m.commandMode
	if m.commandMode {
		m.queryStyle = QueryStyle(25)
	} else {
		m.queryStyle = QueryStyle((m.indecies.QueryIndex % 5) + 10)
	}
	m.queryField.Focus()
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
				case ":":
					return KeyColon(m)
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
				case "up": 
					return KeyDown(m)
				case "down": 
					return KeyUp(m)
				case "tab":
					return KeyTab(m)
				case "ctrl+g":
					return KeyCtrlG(m)
				case "ctrl+f":
					return KeyCtrlF(m)
			}
		}
	m.queryField, cmd = m.queryField.Update(msg)
	return m, cmd	
}

// views

func formViewCategories(m model) string {
	s := strings.Builder{}
	s.WriteString("\n\t1   Select categories:\n\n")
	for i := 0; i < len(coparse.ContextCategories); i++ {
		if m.indecies.FormIndex == i && m.indecies.ContextIndex == 0 {
			s.WriteString("\t[X] ")
		} else if coutils.ContainsInt(m.contextCategories, i) {
			s.WriteString("\t[x] ")	
		} else {
			s.WriteString("\t[ ] ")
		}
		s.WriteString(coparse.ContextCategories[i] + "\n")
	}
	return s.String()
}

func formViewComment(m model) string {
	commentCategories := []string{"include", "don't include"}
	s := strings.Builder{}
	s.WriteString("\n\t2   Include comments:\n\n")
	for i := 0; i < len(commentCategories); i++ {
		if m.indecies.FormIndex == i && m.indecies.ContextIndex == 1 {
			s.WriteString("\t[X] ")
		} else if m.contextComment[i] == 1 {
			s.WriteString("\t[x] ")	
		} else {
			s.WriteString("\t[ ] ")
		}
		s.WriteString(commentCategories[i] + "\n")
	}
	return s.String()
}

func formView(m model) string { 
	s := strings.Builder{}
	s.WriteString(formViewCategories(m))
	s.WriteString(formViewComment(m))
	s.WriteString("\n\tpress ctrl+c to quit | press ctrl+f to return | press enter to submit choice | tab to switch \n")
	return s.String()
}

func searchView(m model, title string) string {
  return lipgloss.Place(
  	m.width,
  	m.height,
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Center, 
			lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				m.queryStyle.InputField.Render(m.queryField.View()),
			),
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.resultStyle.InputField.Render(m.resultField.View()),
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					m.query.ResultLocations[m.indecies.ResultIndex],
					" | ",
					strconv.Itoa(m.indecies.ResultIndex+1),
					"/",
					strconv.Itoa(len(m.query.Result)),
					" | press ctrl+c to quit | press ctrl+f for settings",
				),
			),
		),
	)
}

/* 
** @name: View
** @description: Returns the layout/placement of the visual elements of the TUI.
*/
func (m model) View() string {
	title := m.query.QueryType[m.indecies.QueryIndex]
	if m.commandMode {
		title = "command mode"
	}
	if m.formMode {
		return formView(m) 
	} else {
		return searchView(m, title)
	}
}

// runner function

func main() {
	fmt.Println("Codis (alpha version). Last updated: January 2024 By Timo Kats")
	queryTypes := []string{"Quick search", "Fuzzy search", "Explorative search", "Dependency search", "File view"}
	query := cotypes.Query{"", []string{"None"}, []string{"None"}, queryTypes}
	m := New(query)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

