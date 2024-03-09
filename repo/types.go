package repo

type Repo struct {
	modules map[string]*Module
}

type Module struct {
	Path                      string // module path ie. github.com/perbu/gogrok
	Location                  string // file path
	Version                   string // module version
	Dependencies              []*Module
	Packages                  []*Package
	Type                      DepType
	Repo                      *Repo
	ReverseModuleDependencies []*Module
}

type Package struct {
	Name                string
	Location            string
	Module              *Module
	Files               []*File
	ReverseDependencies []*Package
}

type File struct {
	Name    string
	Imports []*Package
	Package *Package
	Module  *Module
	Lines   []string
}
