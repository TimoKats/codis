package coexplore

import (
	"os"
	"path/filepath"
	"time"
	"strings"
	"strconv"

	coutils "codis/lib/coutils"
	coparse "codis/lib/coparse"
)

// globals

var id = -1
var minLevel = 0
var selectedId = -1
var fileTree = "" 
var selectedPath = ""

// structs

// add ctrl+f to filter this!
type FileInfo struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time   `json:"mod_time"`
	IsDir   bool        `json:"is_dir"`
}

func fileInfoFromInterface(v os.FileInfo) *FileInfo {
	return &FileInfo{v.Name(), v.Size(), v.Mode(), v.ModTime(), v.IsDir()}
}

type Node struct {
	FullPath string    `json:"path"`
	Info     *FileInfo `json:"info"`
	Children []*Node   `json:"children"` 
	Parent   *Node     `json:"-"`
}

// filetree related functions

/* 
** @name: NewTree
** @description: Creates a filetree based on a current directory and a root node object.
*/
func NewTree(root string) (result *Node, err error) { 
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return
	}
	parents := make(map[string]*Node)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		parents[path] = &Node{
			FullPath: path,
			Info:     fileInfoFromInterface(info),
			Children: make([]*Node, 0),
		}
		return nil
	}
	if err = filepath.Walk(absRoot, walkFunc); err != nil {
		return
	}
	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { 
			result = node
		} else {
			node.Parent = parent
			if node.Info.IsDir {
				parent.Children = append(parent.Children, node)
			} else {
				parent.Children = append([]*Node{node}, parent.Children...)
			}
		}
	}
	return
}

/* 
** @name: selectInfoBox 
** @description: Picks and returns the correct infobox as a string.  
*/ // USE FULLPATH THING IN CODEPENDENCIES ALSO, JUST PREPEND THE WORKING DIRECTORY
func selectInfoBox(node *Node, line string, infoIndex int, escape bool) string {
	if escape {
		return coutils.FormatInfoBox(line, "")
	}
	if infoIndex == 0 {
		return coutils.FormatInfoBox(line, coparse.Categories[node.FullPath])
	} else if infoIndex == 1 {
		return coutils.FormatInfoBox(line, strconv.Itoa(coparse.TypeCountsFunction[node.FullPath]))
	} else if infoIndex == 2 {
		return coutils.FormatInfoBox(line, strconv.Itoa(coparse.TypeCountsObject[node.FullPath]))
	} else if infoIndex == 3 {
		return coutils.FormatInfoBox(line, strconv.Itoa(coparse.TypeCountsDomain[node.FullPath]))
	} else {
		return "None"
	}
}

/* 
** @name: printTree
** @description: Creates a string that contains the file tree.
*/
func printTree(node *Node, tablevel string, pastFolders []string, currentLevel int, maxLevel int, dirOnly bool, infoIndex int) {
	id += 1
	line := ""
	idString := strconv.Itoa(id)
	if len(node.Children) == 0 && !dirOnly {
		line = idString + coutils.ResponsiveTab(idString) + "|" + tablevel + "- " + node.Info.Name 
		fileTree += line + selectInfoBox(node, line, infoIndex, false) 
	} else {
		for _, child := range node.Children {
			if !coutils.ContainsString(pastFolders, node.FullPath) {
				line = idString + coutils.ResponsiveTab(idString) + "|" + tablevel + "/ " + node.Info.Name 
				fileTree += line + selectInfoBox(node, line, infoIndex, true) 
				pastFolders = append(pastFolders, node.FullPath)
			}
			if currentLevel < maxLevel && !strings.Contains(node.FullPath, ".git") {
				printTree(child, tablevel+"\t", pastFolders, currentLevel+1, maxLevel, dirOnly, infoIndex)
			}
		}
	}
}

/* 
** @name: selectDirectory
** @description: Picks the correct root directory for the file tree based on a zoom level.
*/
func selectDirectory(selectedLine int, pastFolders []string, node *Node) {
	selectedId += 1
	if selectedId != selectedLine {
		for _, child := range node.Children {
			if !coutils.ContainsString(pastFolders, node.FullPath) {
				pastFolders = append(pastFolders, node.FullPath)
			}
			if !strings.Contains(node.FullPath, ".git") {
				selectDirectory(selectedLine, pastFolders, child)
			}
		}
	} else {
		selectedPath = node.FullPath
	}
}

/* 
** @name: splitFileTree 
** @description: Creates different pages to fit a large filetree (so it returns a slice) 
*/
func splitFileTree(fileTree string, infoIndex int) ([]string, []string) {
	pages := []string{}
	locations := []string{}
	tempPage := coutils.FormatInfoBox("", "showing: " + coparse.InfoBoxCategories[infoIndex])
	for index, line := range strings.Split(fileTree, "\n") {
		tempPage += line + "\n"
		if index % 15 == 0 && index != 0 {
			pages = append(pages, tempPage)
			locations = append(locations, "file explorer")
			tempPage = coutils.FormatInfoBox("", "showing: " + coparse.InfoBoxCategories[infoIndex] + "\n")
		}	
	}
	pages = append(pages, tempPage)
	locations = append(locations, "file explorer")
	return pages, locations
}

// caller function

/* 
** @name: Show
** @description: Caller function that prints the filetree based on some parameters.  
*/
func Show(fullTree *Node, currentLevel int, maxLevel int, zoomLevel string, dirOnly bool, infoIndex int) ([]string, []string) {
	selectedPath = ""
  id, selectedId = -1, -1
  if zoom, err := strconv.Atoi(zoomLevel); err == nil {
		selectDirectory(zoom, []string{}, fullTree)
		var selectedTree, _ =  NewTree(selectedPath)
		fileTree = ""
		id = zoom - 1
    printTree(selectedTree, "\t", []string{}, currentLevel, maxLevel, dirOnly, infoIndex)
  } else {
		fileTree = ""
    printTree(fullTree, "\t", []string{}, currentLevel, maxLevel, dirOnly, infoIndex)
  }
  return splitFileTree(fileTree, infoIndex) 
}

