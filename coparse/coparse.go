package coparse

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

	cotypes "codis/cotypes"
	coutils "codis/coutils"
)

// filters


// objects


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

// check beginning of line
func HasComment(row string) bool {
	commentKeywords := []string{"//", "/*", "* "}
	for _, commentKeyword := range commentKeywords {
		if strings.Contains(row, commentKeyword) {
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

func iterate(path string) ([]string, []string, []string) {
	var texts []string
	var files []string
	var paths []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Print("codis parsing error: ")
			log.Fatalf(err.Error())
		}
		if !info.IsDir() && strings.Contains(info.Name(), ".") && !strings.Contains(path, ".git") { 
			texts = append(texts, readFile(path))
			files = append(files, info.Name())
			paths = append(paths, path)
		}
		return nil
	})
	return texts, files, paths
}

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

func ReturnLabels(directory string) (map[cotypes.RowLabel]string, []cotypes.RowLabel) {
	texts, files, paths := iterate(directory)
	labeledRows, orderedKeys := labelRows(texts, files, paths)
	return labeledRows, orderedKeys
}

func getTopics(text string, fileType string, hasComment bool) []string {
	if (fileType == "code" && hasComment) || fileType == "textual" {
		tokens := coutils.SplitAny(text, " .,{}()[]")
		return coutils.MostCommonTokens(tokens)
	} else {
		return []string{}
	}
}
	
// move to other class maybe. cotopics
func ReturnTopics(labeledRows map[cotypes.RowLabel]string, orderedKeys []cotypes.RowLabel) map[string]string  {
	topics := make(map[string][]string)
	for _, key := range orderedKeys {
		if _, ok := topics[key.FilePath]; ok {
    	topics[key.FilePath] = append(topics[key.FilePath], getTopics(labeledRows[key], key.Category, key.HasComment)...)
    } else {
		  topics[key.FilePath] = getTopics(labeledRows[key], key.Category, key.HasComment)
    }
	}
	return coutils.MapSliceToString(topics)
}
