/* 
** @name: coparse
** @author: Timo Kats
** @description: Has functions related to returning queries.
*/

package cosearch

import (
  "strconv"
  "regexp"
  "strings"
  "math"

  coinit "codis/coinit"
	cotypes "codis/cotypes"
	coutils "codis/coutils"
)

/* 
** @name: formatResult 
** @description: Takes the result (index) and returns a string that shows the lines around it.
*/
func formatResult(index int, labeledRows map[cotypes.RowLabel]string, orderedKeys []cotypes.RowLabel) string {
  result := "\n\n\n"
  for i := index-2; i <= index+2; i++ {
    if i >= 0 && i < len(orderedKeys) { 
      if i == index {
        result += strconv.Itoa(i) + ">  " + coutils.CropString(labeledRows[orderedKeys[i]], 75, "\n")
      } else {
        result += strconv.Itoa(i) + "|  " + coutils.CropString(labeledRows[orderedKeys[i]],75, "\n")
      }
    } 
  } 
  return result
}

/* 
** @name: computeFuzzyScore 
** @description: Computes the highest "fuzzy" score on a line of code given a query. 
** @note: Add a boolean function to shorten the if-statements.
*/
func computeFuzzyScore(line string, query string) int { 
  score, tempScore := 0, 0
  prevIndex, queryIndex := 0, 0
  for index, character := range strings.SplitAfter(line,"") {
    if len(query) - 1 > queryIndex {
      if strings.ToLower(string(query[queryIndex])) == strings.ToLower(character) && math.Abs(float64(prevIndex)-float64(index)) <= 1 {
        score += 1
        queryIndex += 1
        prevIndex = index
      } else if len(query) - 2 > queryIndex && strings.ToLower(string(query[queryIndex + 1])) == strings.ToLower(character) {
        score += 1
        queryIndex += 2
        prevIndex = index
      } else if math.Abs(float64(prevIndex)-float64(index)) >= 2 && strings.ToLower(string(query[0])) == strings.ToLower(character) {
        queryIndex = 1
        score = 1
        prevIndex = index
        if tempScore > score { score = tempScore }
      }
    }
  }
  if tempScore > score { score = tempScore }
  return score
}

/* 
** @name: BasicQuery 
** @description: Returns lines that contain a subquery. 
*/
func BasicQuery(query string) ([]string, []string) {
  var reQuery, err = regexp.Compile(query)
  results, locations := []string{}, []string{}
  if err != nil {
    return []string{"invalid query"}, []string{"None"}
  }
	for index, key := range coinit.OrderedKeys {
	  if reQuery.MatchString(coinit.LabeledRows[key]) {
	    results = append(results, formatResult(index, coinit.LabeledRows, coinit.OrderedKeys))
	    locations = append(locations, key.Filename + ", line " + strconv.Itoa(key.Linenumber))
	  }
	}
	if len(results) == 0 {
	  return []string{"None"}, []string{"None"}
	} else {
	  return results, locations
	}
}

/* 
** @name: FuzzyQuery
** @description: Returns the top (20) lines with the highest fuzzy scores.  
*/
func FuzzyQuery(query string) ([]string, []string) {
  results, locations := []string{}, []string{}
  fuzzyResults := []cotypes.RowLabel{}
  threshold := int(float64(len(query))/2.0)
  for _, key := range coinit.OrderedKeys {
    if computeFuzzyScore(coinit.LabeledRows[key], query) > threshold { 
      fuzzyResults = append(fuzzyResults, key)
    }
  }
  if len(fuzzyResults) == 0 {
	  return []string{"None"}, []string{"None"}
  }
  for _, key := range fuzzyResults {
    index := coutils.FindIndex(coinit.OrderedKeys, key)
    results = append(results, formatResult(index, coinit.LabeledRows, coinit.OrderedKeys))
    locations = append(locations, key.Filename + ", line " + strconv.Itoa(key.Linenumber))
  } 
  return results, locations 
}

