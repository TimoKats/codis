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

	coparse "codis/coparse"
	coutils "codis/coutils"
)

/* 
** @name: formatResult 
** @description: Takes the result (index) and returns a string that shows the lines around it.
*/
func formatResult(index int, labeledRows map[coparse.RowLabel]string, orderedKeys []coparse.RowLabel) string {
  result := "\n\n\n"
  for i := index-1; i <= index+2; i++ {
    if i >= 0 && i < len(orderedKeys) { 
      if i == index {
        result += strconv.Itoa(i) + ">  " + labeledRows[orderedKeys[i]] + "\n" 
      } else {
        result += strconv.Itoa(i) + "|  " + labeledRows[orderedKeys[i]] + "\n" 
      }
    } 
  } 
  return result
}

/* 
** @name: computeFuzzyScore 
** @description: Computes the highest "fuzzy" score on a line of code given a query. 
*/
func computeFuzzyScore(line string, query string) int { // maybe create bool function to prevent long statements
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
func BasicQuery(query string, labeledRows map[coparse.RowLabel]string, orderedKeys []coparse.RowLabel) ([]string, []string) {
  var reQuery, err = regexp.Compile(query)
  results, locations := []string{}, []string{}
  if err != nil {
    return []string{"invalid query"}, []string{"None"}
  }
	for index, key := range orderedKeys {
	  if reQuery.MatchString(labeledRows[key]) {
	    results = append(results, formatResult(index, labeledRows, orderedKeys))
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
func FuzzyQuery(query string, labeledRows map[coparse.RowLabel]string, orderedKeys []coparse.RowLabel) ([]string, []string) {
  results, locations := []string{}, []string{}
  fuzzyResults := []coparse.RowLabel{}
  threshold := int(float64(len(query))/2.0)
  for _, key := range orderedKeys {
    if computeFuzzyScore(labeledRows[key], query) > threshold { // this has to be parameterized
      fuzzyResults = append(fuzzyResults, key)
    }
  }
  if len(fuzzyResults) == 0 {
	  return []string{"None"}, []string{"None"}
  }
  for _, key := range fuzzyResults {
    index := coutils.FindIndex(orderedKeys, key)
    results = append(results, formatResult(index, labeledRows, orderedKeys))
    locations = append(locations, key.Filename + ", line " + strconv.Itoa(key.Linenumber))
  } 
  return results, locations 
}

/* 
** @name: GetFileTypes
** @description: Temporary function 
*/
func GetFileTypes(labeledRows map[coparse.RowLabel]string, orderedKeys []coparse.RowLabel) map[string]int {
  fileTypes := make(map[string]int)
	for _, key := range orderedKeys {
		if _, ok := fileTypes[key.Filetype]; ok {
		  fileTypes[key.Filetype] += 1
    } else {
      fileTypes[key.Filetype] = 1
    }
  }
  return fileTypes
}

/* 
** @name: GetFileCategories
** @description: Temporary function 
*/
func GetFileCategories(labeledRows map[coparse.RowLabel]string, orderedKeys []coparse.RowLabel) map[string]int {
  fileCategories := make(map[string]int)
	for _, key := range orderedKeys {
		if _, ok := fileCategories[key.Category]; ok {
		  fileCategories[key.Category] += 1
    } else {
      fileCategories[key.Category] = 1
    }
  }
  return fileCategories
}
