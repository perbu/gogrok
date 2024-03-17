package analytics

import (
	"fmt"
	"github.com/perbu/gogrok/gitver"
	mf "golang.org/x/mod/modfile"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func New(path string) (*Repo, error) {
	r := &Repo{
		modules:  make(map[string]*Module),
		basePath: path,
	}
	return r, nil
}

func (r *Repo) Parse() error {
	start := time.Now()
	repoDirs, err := os.ReadDir(r.basePath)
	if err != nil {
		return fmt.Errorf("os.ReadDir: %w", err)
	}

	// First pass: parse all go.mod files:
	for _, repoDir := range repoDirs {
		if !repoDir.IsDir() {
			continue
		}
		modulePath := path.Join(r.basePath, repoDir.Name())
		err := r.ParseMod(modulePath)
		if err != nil {
			return fmt.Errorf("r.ParseMod(%s): %w", modulePath, err)
		}
	}
	slog.Info("parsed go.mod files and git metadata", "duration", time.Since(start))
	start = time.Now()
	// Second pass: load and parse all source code:
	for _, name := range r.GetModuleNames() {
		mod, ok := r.GetModule(name)
		if !ok {
			return fmt.Errorf("r.GetModule(%s): not found", name)
		}
		// Load source code for local modules
		if mod.Type == DepTypeLocal {
			err := mod.LoadSource()
			if err != nil {
				return fmt.Errorf("mod.LoadSource: %w", err)
			}
		}
	}
	slog.Info("parsed source code", "duration", time.Since(start))
	// Populate reverse dependencies, both packages and modules.
	start = time.Now()
	r.reverseDeps()
	slog.Info("populated reverse dependencies", "duration", time.Since(start))
	return nil
}

func (r *Repo) ModuleFilter(t DepType, substring string) []Module {
	mods := make([]Module, 0)
	for _, v := range r.modules {
		if v.Type == t && (substring == "" || strings.Contains(v.Path, substring)) {
			mods = append(mods, *v)
		}
	}
	// Sort by module path
	sort.Slice(mods, func(i, j int) bool {
		return mods[i].Path < mods[j].Path
	})
	return mods
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
			versions:                  make([]string, 0),
		}
	}
	m.Type = DepTypeLocal
	m.AddVersion(latestVersion)
	m.Location = modulePath
	for _, require := range file.Require {
		ref, ok := r.modules[require.Mod.Path]
		switch ok {
		case true:
			m.Dependencies = append(m.Dependencies, ref)
			ref.AddVersion(require.Mod.Version)
		case false:
			newDep := &Module{
				Path:                      require.Mod.Path,
				Type:                      DepTypeExternal, // all dependencies are external until proven otherwise
				Repo:                      r,
				ReverseModuleDependencies: make([]*Module, 0),
				versions:                  make([]string, 0),
			}
			newDep.AddVersion(require.Mod.Version)
			m.Dependencies = append(m.Dependencies, newDep)
			r.modules[require.Mod.Path] = newDep
		}
	}
	r.modules[file.Module.Mod.Path] = m
	return nil
}
