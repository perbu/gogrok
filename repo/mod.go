package repo

import (
	"fmt"
	mf "golang.org/x/mod/modfile"
	"os"
	"path/filepath"
)

type DepType int

const (
	DepTypeUnknown DepType = iota
	DepTypeLocal
	DepTypeExternal
	DepTypeStdlib
)

type Module struct {
	Path         string // module path ie. github.com/perbu/gogrok
	Location     string // file path
	Version      string // module version
	Dependencies []*Module
	Packages     []*Package
	Type         DepType
	Repo         *Repo
}

func (r *Repo) GetModule(path string) (*Module, bool) {
	m, ok := r.modules[path]
	if !ok {
		return nil, false
	}
	return m, true
}
func (r *Repo) GetModuleNames() []string {
	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	return names
}

// FindModule finds a module from a full import path.
// It works by checking the module directory. If the module is not found
// it will chop off the last part of the path and try again.
func (r *Repo) FindModule(path string) (string, bool) {
	for {
		m, ok := r.modules[path]
		if ok {
			return m.Path, true
		}
		// Adjusted base case condition to check for filesystem root
		parentPath := filepath.Dir(path)
		if parentPath == path { // Means no further parent, filepath.Dir can't go up
			return "", false
		}
		path = parentPath
	}
}

func (r *Repo) ParseMod(modulePath string) error {
	modFilePath := filepath.Join(modulePath, "go.mod")
	data, err := os.ReadFile(modFilePath)
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}

	file, err := mf.Parse(modulePath, data, nil)
	if err != nil {
		return fmt.Errorf("modfile.Parse: %w", err)
	}
	m, ok := r.modules[file.Module.Mod.Path]
	if !ok {
		m = &Module{
			Path:         file.Module.Mod.Path,
			Location:     modulePath,
			Version:      file.Module.Mod.Version,
			Dependencies: make([]*Module, 0, len(file.Require)),
			Repo:         r,
		}
	}
	m.Type = DepTypeLocal
	for _, require := range file.Require {
		ref, ok := r.modules[require.Mod.Path]
		switch ok {
		case true:
			m.Dependencies = append(m.Dependencies, ref)
		case false:
			newDep := &Module{
				Path:    require.Mod.Path,
				Version: require.Mod.Version,
				Type:    DepTypeExternal, // all dependencies are external until proven otherwise
				Repo:    r,
			}
			m.Dependencies = append(m.Dependencies, newDep)
			r.modules[require.Mod.Path] = newDep
		}
	}
	r.modules[file.Module.Mod.Path] = m
	return nil
}
