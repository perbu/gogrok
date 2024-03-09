package render

type Frontpage struct {
	LocalModules    []Module
	ExternalModules []Module
}

type Module struct {
	Path                string
	Dependencies        []string
	ReverseDependencies []string
}
