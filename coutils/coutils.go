/* 
** @name: coutils 
** @author: Timo Kats
** @description: Functions that are not related to a specific functionality of codis. 
*/

package coutils

import (
  "math"
  "strings"

	coparse "codis/coparse"
)

/* 
** @name: FindIndex 
** @description: Returns the location of a RowLabel index.
*/
func FindIndex(orderedKeys []coparse.RowLabel, queryKey coparse.RowLabel) int {
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
func FindMaxMap(unsortedResults map[coparse.RowLabel]float64) coparse.RowLabel {
  max := 0.0
  var highestKey coparse.RowLabel
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

