/* 
** @name: coutils 
** @author: Timo Kats
** @description: Functions that are not related to a specific functionality of codis. 
*/

package coutils

import (
  "os"
  "bufio"
  "math"
  "strings"
  "unicode"

	cotypes "codis/cotypes"
)

// globals

var Stopwords = readLines("coutils/coimports/stopwords.txt")

func readLines(path string) []string {
    file, err := os.Open(path)
    if err != nil {
        return []string{} 
    }
    defer file.Close()
    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines
}

/* 
** @name: CropString
** @description: Truncnates a string based on a max threshold. 
*/
func CropString(line string, max int, end string) string {
  if len(line) < max {
    return line + end 
  } else {
    return line[:max] + end 
  }
}

func tabCorrectedLen(line string) int {
  len := 0
  for _, char := range line {
    if char == '\t' {
      len += 4
    } else {
      len += 1
    } 
  }
  return len
}

func FormatInfoBox(line string, info string) string {
  len := tabCorrectedLen(line)
  whitespace := strings.Repeat(" ", (60-len))
  return whitespace + "| " + info + "\n"
}

/* 
** @name: FindIndex 
** @description: Returns the location of a RowLabel index.
*/
func FindIndex(orderedKeys []cotypes.RowLabel, queryKey cotypes.RowLabel) int {
  index := 0
  for index, key := range orderedKeys {
    if key == queryKey {
      return index
    } 
  }
  return index
}

/* 
** @name: containsString 
** @description: Returns true if a list contains a string. 
*/
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

/* 
** @name: SplitAny 
** @description: Splits a string based on a set of characters/runes. 
*/
func SplitAny(s string, seps string) []string {
    splitter := func(r rune) bool {
        return strings.ContainsRune(seps, r)
    }
    return strings.FieldsFunc(s, splitter)
}

func hasSymbol(str string) bool {
    for _, letter := range str {
        if unicode.IsSymbol(letter) || letter == '@' || letter == '_' {
            return true
        }
    }
    return false
}

func FormatTopics(mapping map[string]string) map[string]string {
  newMapping := make(map[string]string)
  for key, value := range mapping {
    tokens := SplitAny(value, " _,.;(){}[]")
    tempTokens := []string{}
    for _, token := range tokens {
      if !hasSymbol(token) && !ContainsString(Stopwords, strings.ToLower(token)) {
        tempTokens = append(tempTokens, token)
      }
    }
    newMapping[key] = CropString(strings.Join(tempTokens, ", "), 30, "") 
  }
  return newMapping
}

/* 
** @name: FindMaxSlice 
** @description: Returns the highest value in an unsorted slice. 
*/
func FindMaxSlice(arr []float64) float64 {
   max := 0.0
   for i := 0; i < len(arr); i++ {
      if arr[i] > max && !math.IsNaN(arr[i]) {
         max = arr[i]
      }
   }
   return max
}

/* 
** @name: FindMaxMap 
** @description: Returns the key associated with the highest value in a map. 
*/
func FindMaxMap(unsortedResults map[cotypes.RowLabel]float64) cotypes.RowLabel {
  max := 0.0
  var highestKey cotypes.RowLabel
  for key, score := range unsortedResults {
    if score > max {
      highestKey = key 
      max = score
    }
  }
  return highestKey
}


/* 
** @name: ResponsiveTab 
** @description: Returns whitespace with a corrected offset based on a leading character
*/
func ResponsiveTab(offsetString string) string {
  offset := len(offsetString)
  tabSize := 4
  result := ""
  if offset < tabSize {
    for index := 0; index < tabSize - offset; index++ {
      result += " "
    }
  }
  return result
}

