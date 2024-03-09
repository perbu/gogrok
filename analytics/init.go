package analytics

import (
	"fmt"
	"os"
	"path"
)

func New(path string) (*Repo, error) {
	r := &Repo{
		modules:  make(map[string]*Module),
		basePath: path,
	}
	return r, nil
}

func (r *Repo) Parse() error {
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

	// Second pass: load and parse all source code:
	for _, name := range r.GetModuleNames() {
		mod, ok := r.GetModule(name)
		if !ok {
			return fmt.Errorf("r.GetModule(%s): not found", name)
		}
		if mod.Type == DepTypeLocal {
			err := mod.LoadSource()
			if err != nil {
				return fmt.Errorf("mod.LoadSource: %w", err)
			}
		}
	}
	// Populate reverse dependencies, both packages and modules.
	r.reverseDeps()
	return nil
}
