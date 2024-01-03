package coparse

/* 
** @name: coparse 
** @author: Timo Kats
** @description: Walks through and parses the files in the child-directories. 
*/

import (
	"log"
	"os"
	"fmt"
	"path/filepath"
	"strings"

	coutils "codis/lib/coutils"
	cotypes "codis/lib/cotypes"
)

// globals
var codeStarted = false
var CurrentDirectory, _ = os.Getwd()
var LabeledRows, OrderedKeys = ReturnLabels(CurrentDirectory)
var InfoBoxCategories = []string{"types", "#functions", "#objects", "#web domains"}
var Categories = ReturnCategories(LabeledRows, OrderedKeys)
var TypeCountsFunction = ReturnTypeCounts(LabeledRows, OrderedKeys, "function")
var TypeCountsObject = ReturnTypeCounts(LabeledRows, OrderedKeys, "object")
var TypeCountsDomain = ReturnTypeCounts(LabeledRows, OrderedKeys, "domain")
var Imports = ReturnImports(LabeledRows, OrderedKeys)
var ContextCategories = ReturnUniqueCategories(OrderedKeys)
var InvertedIndex = ReturnIndex()

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
	categories["code"] = []string{"rb","c","cc","cpp","py","js","java","go", "h"}
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
func hasDomain(row string) bool { 
	domainKeywords := []string{"http://", "https://", ".com", ".nl"}
	for _, domainKeyword := range domainKeywords {
		if strings.Contains(row, domainKeyword) {
			return true
		}
	}
	return false
}

// DONE
func getImportedFile(line string, paths []string, files []string) string {
	for index, file := range files { 
		if strings.Contains(line, file) || strings.Contains(line, strings.Split(file, ".")[0]) {
			return paths[index][len(CurrentDirectory):]
		} 
	}
	return ""
}

// DONE
func importedCode(row string, HasComment bool, files []string, paths []string, fileCategory string) string {
	if fileCategory == "code" && !codeStarted && !HasComment && coutils.HasAlpha(row) {
		return getImportedFile(row, paths, files)
	} else {
		return ""
	}
}

/* 
** @name: HasComment 
** @description: Checks if the current line has a comment in it.
*/
func hasComment(row string) bool {
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
func hasVariableDeclaration(row string, fileCategory string) bool { 
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
			codeStarted = true
			return true
		}
	}
	return false
}

/* 
** @name: HasObject 
** @description: Returns true if a line defines an object/class.
*/
func hasObject(row string, fileCategory string) bool { 
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
			codeStarted = true
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
			codeStarted = true
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
		codeStarted = false
		for lineIndex, line := range strings.Split(text, "\n") {
			HasVariableDeclaration := hasVariableDeclaration(line, fileCategory)
			HasObject := hasObject(line, fileCategory)
			hasFunction := hasFunction(line, fileCategory)
			HasDomain := hasDomain(line)
			HasComment := hasComment(line)
			importedCode := importedCode(line, HasComment, files, paths, fileCategory)
			// create key 
			key := cotypes.RowLabel{
				Filename: files[fileIndex], 
				Linenumber: lineIndex+1, 
				Filetype: FiletypeString, 
				HasVariableDeclaration: HasVariableDeclaration,
				HasFunction: hasFunction, 
				HasObject: HasObject, 
				HasDomain: HasDomain, 
				Category: fileCategory,
				HasComment: HasComment, 
				FilePath: FilePath,
				ImportedCode: importedCode,
			}
			// create return values
			labeledRows[key] = line
			orderedKeys = append(orderedKeys, key)
		}
		fmt.Println("Parsed: ", files[fileIndex])
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

func ReturnTypeCounts(labeledRows map[cotypes.RowLabel]string, orderedKeys []cotypes.RowLabel, lineType string) map[string]int {
	counts := make(map[string]int)
	for _, key := range orderedKeys {
		if _, ok := counts[key.FilePath]; ok {
			if lineType == "function" && key.HasFunction {
				counts[key.FilePath] += 1 
			} else if lineType == "object" && key.HasObject {
				counts[key.FilePath] += 1 
			} else if lineType == "domain" && key.HasDomain {
				counts[key.FilePath] += 1 
			}
		} else if lineType == "function" && key.HasFunction {
			counts[key.FilePath] = 1
		} else if lineType == "object" && key.HasObject {
			counts[key.FilePath] = 1
		} else if lineType == "domain" && key.HasDomain {
			counts[key.FilePath] = 1
		} 
	}
	return counts 
}

func ReturnImports(labeledRows map[cotypes.RowLabel]string, orderedKeys []cotypes.RowLabel) map[string][]string {
	imports := make(map[string][]string)
	for _, key := range orderedKeys {
		if key.ImportedCode != "" {
			if _, ok := imports[key.FilePath[len(CurrentDirectory):]]; ok {
				if !coutils.ContainsString(imports[key.FilePath[len(CurrentDirectory):]], key.ImportedCode) {
					imports[key.FilePath[len(CurrentDirectory):]] = append(imports[key.FilePath[len(CurrentDirectory):]], key.ImportedCode)
				}
			} else if !strings.Contains(key.ImportedCode, key.FilePath[len(CurrentDirectory):]) {
				imports[key.FilePath[len(CurrentDirectory):]] = []string{key.ImportedCode}
			}
		}
	}
	return imports 
}

func ReturnUniqueCategories(orderedKeys []cotypes.RowLabel) []string {
	uniqueFileTypes := []string{}
	for _, key := range orderedKeys {
		if !coutils.ContainsString(uniqueFileTypes, key.Category) {
			uniqueFileTypes = append(uniqueFileTypes, key.Category)
		}
	}
	return uniqueFileTypes
}

func ReturnIndex() map[string][]cotypes.IndexLabel {
  invertedIndex := make(map[string][]cotypes.IndexLabel)
  for index, key := range OrderedKeys {
    tokens := coutils.SplitAny(LabeledRows[key], " :;{}().,[]")
    for _, token := range tokens {
    	filename := key.FilePath[len(CurrentDirectory):]
      if _, ok := invertedIndex[token]; !ok && len(token) > 1 {
        invertedIndex[token] = []cotypes.IndexLabel{cotypes.IndexLabel{filename, key.Category, key.Filetype, key.Linenumber, index}}
      } else {
        invertedIndex[token] = append(invertedIndex[token], cotypes.IndexLabel{filename, key.Category, key.Filetype, key.Linenumber, index})
      }
    }
  }
  return invertedIndex
}
