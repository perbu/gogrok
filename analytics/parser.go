package analytics

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func (m *Module) LoadSource() error {
	err := filepath.Walk(m.Location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("error in walkfunc", err)
			fmt.Println("module:", m.Path)
			fmt.Println("location:", m.Location)
			fmt.Println("path:", path)
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// Ensure we know about the package
			p, ok := m.GetPackage(file.Name.Name)
			if !ok {
				rpath, err := filepath.Rel(m.Location, path)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("filepath.Rel: %w", err)
				}
				p = &Package{
					Name:                file.Name.Name,
					Location:            rpath,
					Module:              m,
					Files:               make([]*File, 0),
					ReverseDependencies: make([]*Package, 0),
				}
				m.AddPackage(p)
			}
			// Ensure we know about the file in the package.
			// Every file is new, so we don't need to check for existence.
			f := p.AddFile(path)
			ast.Inspect(file, func(n ast.Node) bool {
				// Find Import Specs
				imp, ok := n.(*ast.ImportSpec)
				if ok {
					f.AddImport(imp.Path.Value)
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
