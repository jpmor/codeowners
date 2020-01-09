package codeowners

type PathSufix int

const (
	Absolute  int = iota // File ends e.g "index.php"
	Flat                 // All other files in directory, "app/lib/*"
	Recursive            // All files in the directory AND all subfiles, "app/lib/"
	Type                 // All files in the path ending in the subsequent file type, "*.rb"
	None                 // Nothing
)
