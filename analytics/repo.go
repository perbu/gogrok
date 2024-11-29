package analytics

import (
	"context"
	"fmt"
	"github.com/perbu/gogrok/modver"
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
		modules:    make(map[string]*Module),
		basePath:   path,
		modTracker: modver.New(),
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
	// Get remove tags for all external dependencies
	start = time.Now()
	for _, mod := range r.ModuleFilter(DepTypeExternal, "") {
		versions, err := r.modTracker.GetTags(context.TODO(), mod.Path)
		if err != nil {
			return fmt.Errorf("modTracker.GetTags(%s): %w", mod.Path, err)
		}
		mod.AddVersions(versions)
		slog.Debug("fetched remote tags", "module", mod.Path, "versions", len(versions))
	}
	slog.Info("fetched remote tags", "duration", time.Since(start))
	// close the modTracker cache:
	err = r.modTracker.Close()
	if err != nil {
		return fmt.Errorf("modTracker.Close: %w", err)
	}
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

	latestVersion, err := r.modTracker.GetLatestVersion(context.TODO(), modulePath)
	if err != nil {
		return fmt.Errorf("gitver.GetLatestTag: %w", err)
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

func (r *Repo) reverseDeps() {
	// Initialize a map for reverse package dependencies to avoid duplicates
	packageReverseDepsMap := make(map[*Package]map[*Package]struct{})

	// Iterate through modules in the repo
	for _, module := range r.modules {
		// Populate Reverse Module Dependencies
		for _, dependency := range module.Dependencies {
			dependency.ReverseModuleDependencies = append(dependency.ReverseModuleDependencies, module)
		}

		// Access packages within the module
		for _, pkg := range module.Packages {
			// Ensure the map for the package is initialized
			if _, exists := packageReverseDepsMap[pkg]; !exists {
				packageReverseDepsMap[pkg] = make(map[*Package]struct{})
			}

			// Populate Reverse Package Dependencies using file imports
			for _, file := range pkg.files {
				for _, importedPkg := range file.Imports {
					// Avoid adding duplicate reverse dependencies
					if _, exists := packageReverseDepsMap[importedPkg][pkg]; !exists {
						if _, exists := packageReverseDepsMap[importedPkg]; !exists {
							packageReverseDepsMap[importedPkg] = make(map[*Package]struct{})
						}
						packageReverseDepsMap[importedPkg][pkg] = struct{}{}
						importedPkg.ReverseDependencies = append(importedPkg.ReverseDependencies, pkg)
					}
				}
			}
		}
	}
}
