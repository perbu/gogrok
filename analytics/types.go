package analytics

type Repo struct {
	modules  map[string]*Module
	basePath string
}

type Module struct {
	Path                      string     // module path ie. github.com/perbu/gogrok
	Location                  string     // file path
	Version                   string     // module version
	Dependencies              []*Module  // list of dependencies
	Packages                  []*Package // list of packages in the module
	Type                      DepType    // either a local (on-disk) or external (remote) module
	Repo                      *Repo      // reference to the repo
	ReverseModuleDependencies []*Module  // List of modules that depend on this module
}

type Package struct {
	Name                string     // package name
	Location            string     // file path, relative to the module
	Module              *Module    // reference to the module
	Files               []*File    // list of files in the package
	ReverseDependencies []*Package // list of packages that depend on this package
}

type File struct {
	Name    string     // file name, not including the path
	Imports []*Package // list of imported packages
	Package *Package   // reference to the package this file belongs to
	Module  *Module    // reference to the module
	Lines   []string   // file contents, split into lines, allow for references to files and lines
}
