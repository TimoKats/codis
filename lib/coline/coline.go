/* 
** @name: coline
** @author: Timo Kats
** @description: Contains functions related to line search 
*/

package coline

import (
  "strings"
  "strconv"
  coparse "codis/lib/coparse"
  coutils "codis/lib/coutils"
)

func formatResult(index int, selectedLine string) string {
  result := "\n\n\n"
  result += selectedLine + "\n\n---\n\n"
  for i := index-2; i <= index+2; i++ {
    if i >= 0 && i < len(coparse.OrderedKeys) { 
      if i == index {
        result += strconv.Itoa(i) + ">  " + coutils.CropString(coparse.LabeledRows[coparse.OrderedKeys[i]], 75, "\n")
      } else {
        result += strconv.Itoa(i) + "|  " + coutils.CropString(coparse.LabeledRows[coparse.OrderedKeys[i]],75, "\n")
      }
    } 
  } 
  return result
}

func getSelectedLine(lineNumber int, FileName string) string {
	for _, key := range coparse.OrderedKeys {
    if key.Filename == FileName && key.Linenumber == lineNumber {
      return coparse.LabeledRows[key]
    }
	}
	return "line not found..."
}

func getRelatedLines(tokens []string, selectedLine string) ([]string, []string) {
  illegalKeywords := []string{"var", "let", "include", "import", "func", "def", "fun"}
  results := []string{}
  locations := []string{}
  for _, token := range tokens {
    if len(token) > 2 && !coutils.ContainsString(illegalKeywords, token) && coutils.HasAlpha(token) {
      for index, key := range coparse.OrderedKeys {
        if key.HasVariableDeclaration || key.HasFunction || key.HasObject {
          if strings.Contains(coparse.LabeledRows[key], token) && !coutils.ContainsString(results, coparse.LabeledRows[key]) {
            results = append(results, formatResult(index, selectedLine))
            locations = append(locations, key.Filename + ", line " + strconv.Itoa(key.Linenumber))
          }
        }
      }
    }
  }
  if len(results) == 1 {
    return []string{"No dependencies found"}, []string{"none"}
  } else {
    return results, locations
  }
}

func SearchLine(fileName string, lineNumber int) ([]string, []string) {
  selectedLine := getSelectedLine(lineNumber, fileName)
  if selectedLine == "line not found..." {
    return []string{selectedLine}, []string{"None"}
  } else {
    selectedLineTokens := coutils.SplitAny(selectedLine, " .:=(){}[]")
    return getRelatedLines(selectedLineTokens, selectedLine)
  }
}

