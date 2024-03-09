package repo

import "strings"

type File struct {
	Name    string
	Imports []*Package
	Package *Package
	Module  *Module
}

// AddFile adds a file to the package
func (p *Package) AddFile(name string) *File {
	for _, file := range p.Files {
		if file.Name == name {
			// should never happen
			panic("file already exists")
		}
	}
	// create a new file:
	f := &File{
		Name:    name,
		Imports: make([]*Package, 0),
		Package: p,
		Module:  p.Module,
	}
	p.Files = append(p.Files, f)
	return f
}

// AddImport adds an import to the file
// it should find the module we're importing from and the package
// and add it to the file's imports
func (f *File) AddImport(name string) {
	for _, pkg := range f.Imports {
		if pkg.Name == name {
			panic("import already exists")
		}
	}

	// find the module name:
	mName, ok := f.Module.Repo.FindModule(name)
	if !ok {
		// The modules we can't locate are either stdlib or external, not our code.
		return
	}
	mod, ok := f.Module.Repo.GetModule(mName)
	if !ok {
		// should never happen
		panic("module not found: " + mName)
	}
	pack := packageNameFromImportPath(name)
	p, ok := mod.GetPackage(pack)
	if !ok {
		// should never happen
		panic("package not found: " + pack)
	}
	f.Imports = append(f.Imports, p)
}

func packageNameFromImportPath(path string) string {
	// package name is the last part of the import path
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
