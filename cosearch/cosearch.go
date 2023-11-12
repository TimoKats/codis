package cosearch

import (
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
