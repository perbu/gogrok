package analytics

import (
	"fmt"
	"github.com/perbu/gogrok/gitver"
	mf "golang.org/x/mod/modfile"
	"os"
	"path/filepath"
)

type DepType int

const (
	DepTypeLocal DepType = iota + 1
	DepTypeExternal
)

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

	latestVersion, err := gitver.GetLatestTag(modulePath)
	if err != nil {
		return fmt.Errorf("gitver.GetLatestTag: %w")
	}

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
			Path:                      file.Module.Mod.Path,
			Location:                  modulePath,
			Dependencies:              make([]*Module, 0, len(file.Require)),
			Repo:                      r,
			ReverseModuleDependencies: make([]*Module, 0),
		}
	}
	m.Type = DepTypeLocal
	m.Version = latestVersion
	for _, require := range file.Require {
		ref, ok := r.modules[require.Mod.Path]
		switch ok {
		case true:
			m.Dependencies = append(m.Dependencies, ref)
		case false:
			newDep := &Module{
				Path:                      require.Mod.Path,
				Version:                   require.Mod.Version,
				Type:                      DepTypeExternal, // all dependencies are external until proven otherwise
				Repo:                      r,
				ReverseModuleDependencies: make([]*Module, 0),
			}
			m.Dependencies = append(m.Dependencies, newDep)
			r.modules[require.Mod.Path] = newDep
		}
	}
	r.modules[file.Module.Mod.Path] = m
	return nil
}

func (m *Module) GetStringDependencies() []string {
	depStrings := make([]string, 0, len(m.Dependencies))
	for _, dep := range m.Dependencies {
		depStrings = append(depStrings, dep.Path)
	}
	return depStrings

}

func (m *Module) GetStringReverseDependencies() []string {
	depStrings := make([]string, 0, len(m.ReverseModuleDependencies))
	for _, dep := range m.ReverseModuleDependencies {
		depStrings = append(depStrings, dep.Path)
	}
	return depStrings
}