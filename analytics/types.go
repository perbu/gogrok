package analytics

import (
	"github.com/perbu/gogrok/modver"
	"go/ast"
)

type Repo struct {
	modules    map[string]*Module
	basePath   string
	modTracker *modver.ModTracker
}

type Module struct {
	Path                      string     // module path ie. github.com/perbu/gogrok
	Location                  string     // file path
	versions                  []string   // module versions
	Dependencies              []*Module  // list of dependencies
	Packages                  []*Package // list of packages in the module
	Type                      DepType    // either a local (on-disk) or external (remote) module
	Repo                      *Repo      // reference to the repo
	ReverseModuleDependencies []*Module  // List of modules that depend on this module
	LatestVersion             string     // latest version of the module
}

type Package struct {
	Name                string     // package name
	Location            string     // file path, relative to the module
	Module              *Module    // reference to the module
	files               []*File    // list of files in the package
	ReverseDependencies []*Package // list of packages that depend on this package
}

type File struct {
	Name    string     // file name, not including the path
	Imports []*Package // list of imported packages
	Package *Package   // reference to the package this file belongs to
	Module  *Module    // reference to the module
	source  []string   // file contents, split into lines, allow for references to files and lines
	ast     *ast.File
	Type    FileType
}

type FileType int

const (
	Unknown FileType = iota
	HumanGo
	GeneratedGo
	TestGo
)

//go:generate stringer -type=FileType

// Overview type provides a high-level overview of the repo
type Overview struct {
	Modules          int
	Packages         int
	Files            int
	LoC              int
	OutdatedModules  int
	SecurityProblems int
}
