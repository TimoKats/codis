package main

import (
	"fmt"
	"os"
	coparse "codis/coparse"
	cosearch "codis/cosearch"
	coview "codis/coview"
)

func printLabels(labeledRows map[coparse.RowLabel]string, orderedKeys []coparse.RowLabel) {
	for _, key := range orderedKeys {
		fmt.Println(key, labeledRows[key])
	}
}

func main() {
	fmt.Println("welcome to codis search systems!")
	currentDirectory, _ := os.Getwd()
	labeledRows, orderedKeys := coparse.ReturnLabels(currentDirectory)

	items := cosearch.GetFileTypes(labeledRows, orderedKeys)
	for key, value := range items {
		fmt.Println(key, value)
	}

	items = cosearch.GetFileCategories(labeledRows, orderedKeys)
	for key, value := range items {
		fmt.Println(key, value)
	}
	coview.Show()
}

