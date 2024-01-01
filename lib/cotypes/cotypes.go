package cotypes

// from coparse

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

