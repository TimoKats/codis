/* 
** @name: cofile
** @author: Timo Kats
** @description: Returns a file view based on a query. 
*/

package cofile

import (
  "strings"

  coutils "codis/lib/coutils"
  coparse "codis/lib/coparse"
)

func Show(query string, index int, contextCategories []string) ([]string, []string) {
  filenames, contents := []string{}, []string{}
  for _, filename := range coparse.OrderedFiles {
		filetype := strings.Split(filename, ".")
		filetypeString := filetype[len(filetype)-1]
		fileCategory := coparse.GetFileCategory(filetypeString)
    if strings.Contains(filename, query) && (len(contextCategories) == 0 || coutils.ContainsString(contextCategories, fileCategory)) {
      filenames = append(filenames, filename)
      contents = append(contents, coparse.FileOverview[filename][index])
    }
  }
  if len(filenames) > 0 {
    return contents, filenames
  } else {
    return []string{"None"}, []string{"None"}
  }
}
