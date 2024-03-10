package render

type ModuleOverviewModule struct {
	Path          string
	PackagesCount int
	FilesCount    int
	LinesOfCode   int
	Packages      []string
	Version       string
}

type ModuleOverview struct {
	Modules []ModuleOverviewModule
}
