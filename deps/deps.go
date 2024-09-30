package deps

import "github.com/edoardottt/depsdev/pkg/depsdev"

type Dependency struct {
	Name           string
	LatestVersion  string
	UsedVersions   []string
	SecurityIssues []string
}

type Manager struct {
	client *depsdev.API
	deps   map[string]Dependency
}

func New() *Manager {
	client := depsdev.NewAPI()
	m := &Manager{
		client: client,
	}
	return m
}

func (m *Manager) AddDependency(name string, versions []string) {
	dep := Dependency{
		Name:         name,
		UsedVersions: versions,
	}
	m.deps[name] = dep
}
