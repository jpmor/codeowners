package codeowners

type PathSufix int

const (
	Absolute  int = iota // File ends
	Flat                 // All other files in directory
	Recursive            //all files in the directory AND all subfiles
	None                 // Nothing
)
