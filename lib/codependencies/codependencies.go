/* 
** @name: codependencies 
** @author: Timo Kats
** @description: Returns an overview of the file dependencies in the project. 
*/

package codependencies

import (
  "strings"
  "strconv"

  coparse "codis/lib/coparse"
  coutils "codis/lib/coutils"
)

// globals

var dependencyTree string
var antiCircularDependencies []string
var id = -1

func selectInfoBox(filepath string, line string, infoIndex int) string {
	if infoIndex == 0 {
		return coutils.FormatInfoBox(line, coparse.Categories[coparse.CurrentDirectory + filepath])
	} else if infoIndex == 1 {
		return coutils.FormatInfoBox(line, strconv.Itoa(coparse.TypeCountsFunction[coparse.CurrentDirectory + filepath]))
	} else if infoIndex == 2 {
		return coutils.FormatInfoBox(line, strconv.Itoa(coparse.TypeCountsObject[coparse.CurrentDirectory + filepath]))
	} else if infoIndex == 3 {
		return coutils.FormatInfoBox(line, strconv.Itoa(coparse.TypeCountsDomain[coparse.CurrentDirectory + filepath]))
	} else if infoIndex == 4 {
		return coutils.FormatInfoBox(line, strconv.Itoa(coparse.QueryCounts[coparse.CurrentDirectory + filepath]))
	} else {
		return "None"
	}
}

func GetRootFiles() []string {
	rootFiles := []string{}
	for rootFile, _ := range coparse.Imports {
		unique := true
		for _, importedFiles := range coparse.Imports {
			for _, importedFile := range importedFiles {
				if rootFile == importedFile {
					unique = false
				}
			}
		}
		if unique { rootFiles = append(rootFiles, rootFile) }
	}
	return rootFiles
}

func formatImports(rootFile string, tabLevel string, infoIndex int) {
	antiCircularDependencies = append(antiCircularDependencies, rootFile)
	for _, file := range coparse.Imports[rootFile] {
	  line := ""
		if _, ok := coparse.Imports[file]; ok {
	    id += 1
	    idString := strconv.Itoa(id)
	    line = idString + coutils.ResponsiveTab(idString) + "|" + tabLevel + file
		  dependencyTree += line + selectInfoBox(file, line, infoIndex)
		  if !coutils.ContainsString(antiCircularDependencies, file) {
				formatImports(file, tabLevel + "\t", infoIndex)
		  }
		} else { // it's a leaf
	    id += 1
	    idString := strconv.Itoa(id)
	    line = idString + coutils.ResponsiveTab(idString) + "|" + tabLevel + file
		  dependencyTree += line + selectInfoBox(file, line, infoIndex)
			formatImports(file, tabLevel + "\t", infoIndex)
		}
	}
}

func splitDependencyTree(dependencyTree string, infoIndex int) ([]string, []string) {
	pages := []string{}
	locations := []string{}
	tempPage := coutils.FormatInfoBox("", "showing: " + coparse.InfoBoxCategories[infoIndex])
	for index, line := range strings.Split(dependencyTree, "\n") {
		tempPage += line + "\n"
		if index % 15 == 0 && index != 0 {
			pages = append(pages, tempPage)
			locations = append(locations, "dependency explorer")
			tempPage = coutils.FormatInfoBox("", "showing: " + coparse.InfoBoxCategories[infoIndex] + "\n")
		}	
	}
	pages = append(pages, tempPage)
	locations = append(locations, "dependency explorer")
	return pages, locations
}

func queryRootFile(rootFiles []string, query string) string {
	for _, rootFile := range rootFiles {
		if strings.Contains(rootFile, query) && len(query) > 0 {
			return rootFile
		}
	}
	return ""
}

func Show(infoIndex int, rootFiles []string, query string) ([]string, []string) {
  id = 0
	dependencyTree = ""
	idString := strconv.Itoa(id)
	queriedRootFile := queryRootFile(rootFiles, query)
	antiCircularDependencies = []string{}
	if queriedRootFile == "" && len(rootFiles) > 0 {
		for _, rootFile := range rootFiles {
			antiCircularDependencies = []string{}
			idString = strconv.Itoa(id)
			dependencyTree += idString + coutils.ResponsiveTab(idString) + "|> " + rootFile + "\n"
			formatImports(rootFile, "\t", infoIndex)
			id += 1
		}
	} else {
		dependencyTree += idString + coutils.ResponsiveTab(idString) + "|> " + queriedRootFile + "\n"
		formatImports(queriedRootFile, "\t", infoIndex)
	}
	return splitDependencyTree(dependencyTree, infoIndex)
}
