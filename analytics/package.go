package analytics

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
