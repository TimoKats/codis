package coinit

import (
  "os"
	coparse "codis/coparse"
)

var CurrentDirectory, _ = os.Getwd()
var LabeledRows, OrderedKeys = coparse.ReturnLabels(CurrentDirectory)
var Topics = coparse.ReturnTopics(LabeledRows, OrderedKeys)
var Categories = coparse.ReturnCategories(LabeledRows, OrderedKeys)
