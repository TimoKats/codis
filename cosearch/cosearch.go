package cosearch

import (
  "strings"
  "strconv"

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

func formatResult(index int, labeledRows map[coparse.RowLabel]string, orderedKeys []coparse.RowLabel) string {
  result := ""
  result = strconv.Itoa(index-1) + ": " + labeledRows[orderedKeys[index-1]] + "\n" + strconv.Itoa(index) + ": " + labeledRows[orderedKeys[index]] + "\n" + strconv.Itoa(index+1) + ": " + labeledRows[orderedKeys[index+1]]
  return result
}

func BasicQuery(query string, labeledRows map[coparse.RowLabel]string, orderedKeys []coparse.RowLabel) ([]string, []string) {
  results := []string{}
  locations := []string{}
	for index, key := range orderedKeys {
	  if strings.Contains(labeledRows[key], query) {
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
