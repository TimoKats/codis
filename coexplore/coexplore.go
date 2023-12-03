package coexplore

import (
	"os"
	"path/filepath"
	"time"
	"strings"
	"strconv"
	coutils "codis/coutils"
)

// globals

var id = -1
var minLevel = 0
var selectedId = -1
var fileTree = "" 
var selectedPath = ""

// structs

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

// make a boolean funcion that check the level-thing
func printTree(node *Node, tablevel string, pastFolders []string, currentLevel int, maxLevel int, dirOnly bool) {
	id += 1
	idString := strconv.Itoa(id)
	if len(node.Children) == 0 && !dirOnly {
		fileTree += idString + coutils.ResponsiveTab(idString) + "|" + tablevel + "- " + node.Info.Name + "\n"
	} else {
		for _, child := range node.Children {
			if !coutils.ContainsString(pastFolders, node.FullPath) {
				fileTree += idString + coutils.ResponsiveTab(idString) + "|" + tablevel + "/ " + node.Info.Name + "\n"
				pastFolders = append(pastFolders, node.FullPath)
			}
			if currentLevel < maxLevel && !strings.Contains(node.FullPath, ".git") {
				printTree(child, tablevel+"\t", pastFolders, currentLevel+1, maxLevel, dirOnly)
			}
		}
	}
}

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

func splitFileTree(fileTree string) ([]string, []string) {
	pages := []string{}
	locations := []string{}
	tempPage := ""
	for index, line := range strings.Split(fileTree, "\n") {
		tempPage += line + "\n"
		if index % 15 == 0 && index != 0 {
			pages = append(pages, tempPage)
			locations = append(locations, "file explorer")
			tempPage = ""
		}	
	}
	pages = append(pages, tempPage)
	locations = append(locations, "file explorer")
	return pages, locations
}

// generate fulltree one and pass as parameter!
func Show(fullTree *Node, currentLevel int, maxLevel int, zoomLevel string, dirOnly bool) ([]string, []string) {
	selectedPath = ""
  id, selectedId = -1, -1
  if zoom, err := strconv.Atoi(zoomLevel); err == nil {
		selectDirectory(zoom, []string{}, fullTree)
		var selectedTree, _ =  NewTree(selectedPath)
		fileTree = ""
		id = zoom - 1
    printTree(selectedTree, "\t", []string{}, currentLevel, maxLevel, dirOnly)
  } else {
		fileTree = ""
    printTree(fullTree, "\t", []string{}, currentLevel, maxLevel, dirOnly)
  }
  return splitFileTree(fileTree) 
}

