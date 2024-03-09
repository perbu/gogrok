package main

import (
	"context"
	"fmt"
	"github.com/perbu/gogrok/render"
	"github.com/perbu/gogrok/repo"
	"io"
	"os"
	"os/signal"
	"path"
	"sort"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := run(ctx, os.Stdout, os.Environ(), os.Args)
	if err != nil {
		fmt.Println("run() returned an error: ", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, output io.Writer, env, args []string) error {
	r, err := repo.New()
	if err != nil {
		return fmt.Errorf("repoDir.New: %w", err)
	}
	const codePath = "code"
	repoDirs, err := os.ReadDir(codePath)
	if err != nil {
		return fmt.Errorf("os.ReadDir: %w", err)
	}

	// First pass: parse all go.mod files:
	for _, repoDir := range repoDirs {
		if !repoDir.IsDir() {
			continue
		}
		modulePath := path.Join(codePath, repoDir.Name())
		err := r.ParseMod(modulePath)
		if err != nil {
			return fmt.Errorf("r.ParseMod(%s): %w", modulePath, err)
		}
	}

	localModulesNames := make([]string, 0)
	// Second pass: load and parse all source code:
	for _, name := range r.GetModuleNames() {
		mod, ok := r.GetModule(name)
		if !ok {
			return fmt.Errorf("r.GetModule(%s): not found", name)
		}
		if mod.Type == repo.DepTypeLocal {
			err := mod.LoadSource()
			if err != nil {
				return fmt.Errorf("mod.LoadSource: %w", err)
			}
			localModulesNames = append(localModulesNames, name)
		}
	}
	// sort localModulesNames
	sort.Strings(localModulesNames)

	fpModules := make([]*repo.Module, 0)
	for _, mod := range localModulesNames {
		m, ok := r.GetModule(mod)
		if !ok {
			return fmt.Errorf("r.GetModule(%s): not found", mod)
		}
		fpModules = append(fpModules, m)
	}
	const basePath = "build"
	// Now we can generate the output.
	fp := render.Frontpage{
		Modules: fpModules,
	}
	err = render.RenderFrontpage(fp, path.Join(basePath, "index.html"))
	if err != nil {
		return fmt.Errorf("render.RenderFrontpage: %w", err)
	}

	return nil
}
