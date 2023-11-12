package coparse

// add properties in the RowLabeler, like hasDeclaration, hasClassname, hasFunctionname, hasParameters, hasColumnheader, etc... think about this!
// you can add this in the rowlabeler and then add the property to struct. That should work. KTHXBAI!
// TIP: Work with indentation and {}s! Because that can be indicative of classes/funcitons/etc.
// TIP: Check imports in code files, URLs in all files are ok to search

import (
	"log"
	"os"
	"fmt"
	"path/filepath"
	"strings"
)

// filters


// objects

type RowLabel struct {
	Filename    string
	Filetype    string
	Category    string
	HasObject   bool
	HasDomain   bool
	hasFunction bool
	HasVariableDeclaration bool
	Linenumber  int
}

// i/o functions

func readFile(Filename string) string {
	fileContent, err := os.ReadFile(Filename)
	if err != nil {
		log.Fatal(err)
	}
	text := string(fileContent)
	return text
}

// labeling functions

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

func HasDomain(row string) bool { 
	domainKeywords := []string{"http://", "https://"}
	for _, domainKeyword := range domainKeywords {
		if strings.Contains(row, domainKeyword) {
			return true
		}
	}
	return false
}

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

func iterate(path string) ([]string, []string) {
	var texts []string
	var files []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Print("codis parsing error: ")
			log.Fatalf(err.Error())
		}
		if !info.IsDir() && strings.Contains(info.Name(), ".") { 
			texts = append(texts, readFile(path))
			files = append(files, info.Name())
		}
		return nil
	})
	return texts, files
}

func labelRows(texts []string, files []string) (map[RowLabel]string, []RowLabel) {
	labeledRows := make(map[RowLabel]string)
	orderedKeys := []RowLabel{}
	for fileIndex, text := range texts {
		Filetype := strings.Split(files[fileIndex], ".")
		FiletypeString := Filetype[len(Filetype)-1]
		fileCategory := getFileCategory(FiletypeString)
		for lineIndex, line := range strings.Split(text, "\n") {
			HasVariableDeclaration := HasVariableDeclaration(line, fileCategory)
			HasObject := HasObject(line, fileCategory)
			hasFunction := hasFunction(line, fileCategory)
			HasDomain := HasDomain(line)
			// create key 
			key := RowLabel{Filename: files[fileIndex], 
			Linenumber: lineIndex, Filetype: FiletypeString, 
			HasVariableDeclaration: HasVariableDeclaration,
			hasFunction: hasFunction, HasObject: HasObject, 
			HasDomain: HasDomain, Category: fileCategory}
			// create return values
			labeledRows[key] = line
			orderedKeys = append(orderedKeys, key)
		}
	}
	return labeledRows, orderedKeys
}

// caller functions

func ReturnLabels(directory string) (map[RowLabel]string, []RowLabel) {
	texts, files := iterate(directory)
	labeledRows, orderedKeys := labelRows(texts, files)
	return labeledRows, orderedKeys
}
