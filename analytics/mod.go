package analytics

import (
	"slices"
)

type DepType int

const (
	DepTypeLocal DepType = iota + 1
	DepTypeExternal
)

func (m *Module) Lines() int {
	lines := 0
	for _, pkg := range m.Packages {
		lines += pkg.Lines()
	}
	return lines
}

func (m *Module) Files() int {
	files := 0
	for _, pkg := range m.Packages {
		files += pkg.Files()
	}
	return files
}

func (m *Module) GetPackage(name string) (*Package, bool) {
	for _, pkg := range m.Packages {
		if pkg.Name == name {
			return pkg, true
		}
	}
	return nil, false
}

func (m *Module) AddPackage(p *Package) {
	for _, pkg := range m.Packages {
		if pkg.Name == p.Name {
			return
		}
	}
	m.Packages = append(m.Packages, p)
}

func (m *Module) AddVersion(version string) {
	for _, v := range m.versions {
		if v == version {
			return
		}
	}
	m.versions = append(m.versions, version)
	// sort the versions
	slices.Sort(m.versions)
}

func (m *Module) GetVersions() []string {
	return m.versions
}

func (m *Module) Latest() string {
	if len(m.versions) > 0 {
		return m.versions[len(m.versions)-1]
	}
	return ""
}

func (m *Module) CalculateComplexity() float32 {
	complexity := float32(0)
	for _, pkg := range m.Packages {
		complexity += pkg.CalculateComplexity()
	}
	return complexity / float32(len(m.Packages))
}
