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
