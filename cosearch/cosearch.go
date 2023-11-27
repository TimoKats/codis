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
func computeFuzzyScore(line string, query string) float64 {
  scores := []float64{}
  for _, word := range coutils.SplitAny(line, " ,.()[]{}/") {
    score := 0.0
    for _, character := range strings.SplitAfter(query,"") {
      if strings.Contains(strings.ToLower(word), strings.ToLower(character)) && len(word) < len(query) * 4 {
        score += 1.0
      }
    }
    scores = append(scores, score)
  }
  return coutils.FindMaxSlice(scores)
}

/* 
** @name: sortFuzzyResults
** @description: returns the top (20) results from a collection of fuzzy scores. 
*/
func sortFuzzyResults(unsortedResults map[coparse.RowLabel]float64) []coparse.RowLabel {
  var topResults []coparse.RowLabel
  for i := 0; i < 20; i++ {
    maxKey := coutils.FindMaxMap(unsortedResults)
    topResults = append(topResults, maxKey)
    delete(unsortedResults, maxKey)
  } 
  return topResults
}


/* 
** @name: BasicQuery 
** @description: Returns lines that contain a subquery. 
*/
func BasicQuery(query string, labeledRows map[coparse.RowLabel]string, orderedKeys []coparse.RowLabel) ([]string, []string) {
  var reQuery, err = regexp.Compile(query)
  if err != nil {
    return []string{"invalid query"}, []string{"None"}
  }
  results := []string{}
  locations := []string{}
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
  results := []string{}
  locations := []string{}
  tempResults := make(map[coparse.RowLabel]float64)
  for _, key := range orderedKeys {
    tempResults[key] = computeFuzzyScore(labeledRows[key], query)
  }
  sortedResults := sortFuzzyResults(tempResults) 
  for _, key := range sortedResults {
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
		  fileTypes[key.Filetype]++
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
		  fileCategories[key.Category]++
    } else {
      fileCategories[key.Category] = 1
    }
  }
  return fileCategories
}
