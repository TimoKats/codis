package cosearch

import (
  "strconv"
  "regexp"
  "strings"
  "math"

	coparse "codis/coparse"
)

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

// helper functions

func findIndex(orderedKeys []coparse.RowLabel, queryKey coparse.RowLabel) int {
  index := 0
  for index, key := range orderedKeys {
    if key == queryKey {
      return index
    } 
  }
  return index
}

func splitAny(s string, seps string) []string {
    splitter := func(r rune) bool {
        return strings.ContainsRune(seps, r)
    }
    return strings.FieldsFunc(s, splitter)
}

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

func findMax(arr []float64) float64 {
   max := 0.0
   for i := 0; i < len(arr); i++ {
      if arr[i] > max && !math.IsNaN(arr[i]) {
         max = arr[i]
      }
   }
   return max
}

func findMaxMap(unsortedResults map[coparse.RowLabel]float64) coparse.RowLabel {
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

func computeFuzzyScore(line string, query string) float64 {
  scores := []float64{}
  for _, word := range splitAny(line, " ,.()[]{}/") {
    score := 0.0
    for _, character := range strings.SplitAfter(query,"") {
      if strings.Contains(strings.ToLower(word), strings.ToLower(character)) && len(word) < len(query)*5 {
        score += 1.0
      }
    }
    scores = append(scores, score)
  }
  return findMax(scores)
}

func sortFuzzyResults(unsortedResults map[coparse.RowLabel]float64) []coparse.RowLabel {
  var topResults []coparse.RowLabel
  for i := 0; i < 20; i++ {
    maxKey := findMaxMap(unsortedResults)
    topResults = append(topResults, maxKey)
    delete(unsortedResults, maxKey)
  } 
  return topResults
}

// caller functions

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

func FuzzyQuery(query string, labeledRows map[coparse.RowLabel]string, orderedKeys []coparse.RowLabel) ([]string, []string) {
  results := []string{}
  locations := []string{}
  tempResults := make(map[coparse.RowLabel]float64)
  for _, key := range orderedKeys {
    tempResults[key] = computeFuzzyScore(labeledRows[key], query)
  }
  sortedResults := sortFuzzyResults(tempResults) 
  for _, key := range sortedResults {
    index := findIndex(orderedKeys, key)
    results = append(results, formatResult(index, labeledRows, orderedKeys))
    locations = append(locations, key.Filename + ", line " + strconv.Itoa(key.Linenumber))
  } 
  return results, locations 
}
