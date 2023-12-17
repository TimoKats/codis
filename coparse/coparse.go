package coparse

/* 
** @name: coparse 
** @author: Timo Kats
** @description: Walks through and parses the files in the child-directories. 
*/

// add properties in the cotypes.RowLabeler, like hasDeclaration, hasClassname, hasFunctionname, hasParameters, hasColumnheader, etc... think about this!
// you can add this in the rowlabeler and then add the property to struct. That should work. KTHXBAI!
// TIP: Work with indentation and {}s! Because that can be indicative of classes/funcitons/etc.
// TIP: Check imports in code files, URLs in all files are ok to search

import (
	"log"
	"os"
	"fmt"
	"path/filepath"
	"strings"

	coutils "codis/coutils"
	cotypes "codis/cotypes"
)

// i/o functions

/* 
** @name: readFile
** @description: Returns a string of the filecontents given its name. 
*/
func readFile(Filename string) string {
	fileContent, err := os.ReadFile(Filename)
	if err != nil {
		log.Fatal(err)
	}
	text := string(fileContent)
	return text
}

// labeling functions

/* 
** @name: getFileCategory
** @description: Assigns a category to a file based on its extension 
*/
func getFileCategory(currentExtension string) string {
	var categories = make(map[string][]string)
	categories["data"] = []string{"csv","json","sql","xml","xls","xlsx"}
	categories["web"] = []string{"html","css","scss", "erb"}
	categories["code"] = []string{"rb","c","cc","cpp","py","js","java","go"}
	categories["compiled"] = []string{"dll","exe"}
	categories["textual"] = []string{"txt","md"}
	for Category, extensions := range categories {
		for  _,extension := range extensions {
			if extension == currentExtension {
				return Category
			}
		}
	}
	return "undefined"
}

/* 
** @name: HasDomain
** @description: Returns true if a line has a webdomain. 
*/
func HasDomain(row string) bool { 
	domainKeywords := []string{"http://", "https://"}
	for _, domainKeyword := range domainKeywords {
		if strings.Contains(row, domainKeyword) {
			return true
		}
	}
	return false
}

// check beginning of line
/* 
** @name: HasComment 
** @description: Checks if the current line has a comment in it.
*/
func HasComment(row string) bool {
	commentKeywords := []string{"//", "/*", "* "}
	for _, commentKeyword := range commentKeywords {
		if strings.Contains(row, commentKeyword) {
			return true
		}
	}
	return false
}

/* 
** @name: HasVariableDeclaration
** @description: Returns true if a line has a variable declaration
*/
func HasVariableDeclaration(row string, fileCategory string) bool { 
	declaritiveKeywords := []string{":=","=","let ","var "}
	illegalKeywords := []string{"==", "!=", "//", "/*", "#"}
	if fileCategory != "code" {
		return false
	}
	for _, illegalKeyword := range illegalKeywords {
		if strings.Contains(row, illegalKeyword) {
			return false
		}
	}
	for _, declaritiveKeyword := range declaritiveKeywords {
		if strings.Contains(row, declaritiveKeyword) {
			return true
		}
	}
	return false
}

/* 
** @name: HasObject 
** @description: Returns true if a line defines an object/class.
*/
func HasObject(row string, fileCategory string) bool { 
	declaritiveKeywords := []string{"class ", "struct ", "enum "}
	illegalKeywords := []string{"=", ":=", "//", "/*", "#"}
	if fileCategory != "code" {
		return false
	}
	for _, illegalKeyword := range illegalKeywords {
		if strings.Contains(row, illegalKeyword) {
			return false
		}
	}
	for _, declaritiveKeyword := range declaritiveKeywords {
		if strings.Contains(row, declaritiveKeyword) {
			return true
		}
	}
	return false
}

/* 
** @name: hasFunction
** @description: Returns true if the current line declares a function. 
*/
func hasFunction(row string, fileCategory string) bool { 
	declaritiveKeywords := []string{"def ", "fun ", "fn ", "func "}
	illegalKeywords := []string{"//", "/*", "#"}
	if fileCategory != "code" {
		return false
	}
	for _, illegalKeyword := range illegalKeywords {
		if strings.Contains(row, illegalKeyword) {
			return false
		}
	}
	for _, declaritiveKeyword := range declaritiveKeywords {
		if strings.Contains(row, declaritiveKeyword) {
			return true
		}
	}
	return false
}

// labeler functions

/* 
** @name: iterate
** @description: Walks through all files in the child directories and returns their contents
*/
func iterate(path string) ([]string, []string, []string) {
	var texts []string
	var files []string
	var paths []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Print("codis parsing error: ")
			log.Fatalf(err.Error())
		}
		if !info.IsDir() && strings.Contains(info.Name(), ".") && !strings.Contains(path, ".exe") && !strings.Contains(path, ".git") { 
			texts = append(texts, readFile(path))
			files = append(files, info.Name())
			paths = append(paths, path)
		}
		return nil
	})
	return texts, files, paths
}

/* 
** @name: labelRows
** @description: Creates a rowlabel object for each line of read content.
*/
func labelRows(texts []string, files []string, paths []string) (map[cotypes.RowLabel]string, []cotypes.RowLabel) {
	labeledRows := make(map[cotypes.RowLabel]string)
	orderedKeys := []cotypes.RowLabel{}
	for fileIndex, text := range texts {
		// attributes that are the same for all files
		Filetype := strings.Split(files[fileIndex], ".")
		FiletypeString := Filetype[len(Filetype)-1]
		fileCategory := getFileCategory(FiletypeString)
		FilePath := paths[fileIndex]
		for lineIndex, line := range strings.Split(text, "\n") {
			HasVariableDeclaration := HasVariableDeclaration(line, fileCategory)
			HasObject := HasObject(line, fileCategory)
			hasFunction := hasFunction(line, fileCategory)
			HasDomain := HasDomain(line)
			HasComment := HasComment(line)
			// create key 
			key := cotypes.RowLabel{Filename: files[fileIndex], 
			Linenumber: lineIndex, Filetype: FiletypeString, 
			HasVariableDeclaration: HasVariableDeclaration,
			HasFunction: hasFunction, HasObject: HasObject, 
			HasDomain: HasDomain, Category: fileCategory,
			HasComment: HasComment, FilePath: FilePath}
			// create return values
			labeledRows[key] = line
			orderedKeys = append(orderedKeys, key)
		}
	}
	return labeledRows, orderedKeys
}

// caller functions

/* 
** @name: ReturnLabels
** @description: Returns the rowlabel objects and their order.
*/
func ReturnLabels(directory string) (map[cotypes.RowLabel]string, []cotypes.RowLabel) {
	texts, files, paths := iterate(directory)
	labeledRows, orderedKeys := labelRows(texts, files, paths)
	return labeledRows, orderedKeys
}

/* 
** @name: rankLine
** @description: Assigns a rank to a line that resembles its informative content.
*/
func rankLine(line string, Linenumber int, HasComment bool, category string) (string, int) {
	rank := 0
	if category == "code" {
		if HasComment {
			rank += 1
		}
		if Linenumber < 10 {
			rank += 2
		}
		if strings.Contains(line, "description") && Linenumber < 10 {
			rank += 5
		}
		return line, rank
	} else {
		return "None", rank
	}
}
	
/* 
** @name: ReturnTopics
** @description: Returns the topics of a file 
*/
func ReturnTopics(labeledRows map[cotypes.RowLabel]string, orderedKeys []cotypes.RowLabel) map[string]string  {
	topics := make(map[string]string)
	ranks := make(map[string]int)
	for _, key := range orderedKeys {
		line, rank := rankLine(labeledRows[key], key.Linenumber, key.HasComment, key.Category)
		if _, ok := topics[key.FilePath]; ok { 
    	if rank > ranks[key.FilePath] { 
    		topics[key.FilePath] = line
    		ranks[key.FilePath] = rank
    	} 
    } else if rank >= 1 {
    	topics[key.FilePath] = line
		  ranks[key.FilePath] = rank 
    }
	}
	return coutils.FormatTopics(topics)
}

/* 
** @name: ReturnCategories
** @description: Returns a map with the topic for each filepath 
*/
func ReturnCategories(labeledRows map[cotypes.RowLabel]string, orderedKeys []cotypes.RowLabel) map[string]string  {
	categories := make(map[string]string)
	for _, key := range orderedKeys {
		if _, ok := categories[key.FilePath]; !ok {
			categories[key.FilePath] = key.Category
    } 
	}
	return categories 
}

// init parser (functions can be private now...)

var CurrentDirectory, _ = os.Getwd()
var LabeledRows, OrderedKeys = ReturnLabels(CurrentDirectory)
var Topics = ReturnTopics(LabeledRows, OrderedKeys)
var Categories = ReturnCategories(LabeledRows, OrderedKeys)

