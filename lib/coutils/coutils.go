/* 
** @name: coutils 
** @author: Timo Kats
** @description: Functions that are not related to a specific functionality of codis. 
*/

package coutils

import (
  "math"
  "strings"
  "unicode"

	cotypes "codis/lib/cotypes"
)

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
** @name: containsint
** @description: Returns true if a list contains an integer.
*/
func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func DeleteInt(s []int, e int) []int {
  newSlice := []int{}
  for _, item := range s {
    if item != e {
      newSlice = append(newSlice, item)
    }
  }
  return newSlice
}

func SubsetSlice(s []string, i []int) []string {
  newSlice := []string{}
  for _, index := range i {
    newSlice = append(newSlice, s[index])
  }
  return newSlice
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

func HasAlpha(str string) bool {
    for _, letter := range str {
        if !unicode.IsSymbol(letter) {
            return true
        }
    }
    return false
}

func hasSymbol(str string) bool {
    for _, letter := range str {
        if unicode.IsSymbol(letter) || letter == '@' || letter == '_' {
            return true
        }
    }
    return false
}

func GetJSONFieldname(line string) string {
  if strings.Contains(line, ":") {
    return strings.Split(line, ":")[0]
  } else {
    return "" 
  }
}

func GetCSVSeperator(line string) string {
  seperatorCandidates := map[string]int{",":0,";":0,"\t":0}
  for _, letter := range line {
		if _, ok := seperatorCandidates[string(letter)]; ok {
		  seperatorCandidates[string(letter)] += 1
		}
  }
  return FindMaxMapInt(seperatorCandidates)
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

func FindMaxMapInt(unsortedResults map[string]int) string {
  max := 0
  var highestKey string 
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

