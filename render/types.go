package render

type ModuleOverviewModule struct {
	Path          string
	PackagesCount int
	FilesCount    int
	LinesOfCode   int
	Packages      []string
	Version       string
	Deps          []string
	RevDeps       []string
}

type ModuleOverview struct {
	Modules []ModuleOverviewModule
}

type ModuleDetailModule struct {
	Path          string
	Version       string
	Packages      []ModuleDetailModulePackage
	PackagesCount int
	FilesCount    int
	LinesOfCode   int
}

type ModuleDetailModulePackage struct {
	Name   string
	Module string
}

type Package struct {
	Name     string // The name of the package
	Location string // The path of the package, within the module
	Module   string
	Files    []PackageFile // A slice of files contained within the package
}

type PackageFile struct {
	Name    string
	Package string
	Module  string
}
