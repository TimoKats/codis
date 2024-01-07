/* 
** @name: cotypes
** @author: Timo Kats
** @description: All the objects that are used in different files of codis. 
*/

package cotypes

type RowLabel struct {
	Filename    string
	Filetype    string
	FilePath 		string
	Category    string
	HasObject   bool
	HasDomain   bool
	HasFunction bool
	HasComment 	bool
	ImportedCode 	string
	HasVariableDeclaration bool
	Linenumber  int
}

type Indecies struct {
	QueryIndex int
	ResultIndex int
	FormIndex int
	InfoIndex int
	ContextIndex int
	FileViewIndex int
}

type Query struct {
	Query string
	Result []string
	ResultLocations []string
	QueryType []string
}

/* index not for this version
type IndexLabel struct {
  Filename    string
  Category 		string
  Filetype 		string
  Linenumber  int
  Index       int
}
*/
