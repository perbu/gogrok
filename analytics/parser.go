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
		if info == nil {
			return nil
		}
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		m.NoOfFiles++
		fset := token.NewFileSet()
		astFile, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Println(err)
			return err
		}
		// Ensure we know about the package
		p, ok := m.GetPackage(astFile.Name.Name)
		if !ok {
			rpath, err := filepath.Rel(m.Location, path)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("filepath.Rel: %w", err)
			}
			p = &Package{
				Name:                astFile.Name.Name,
				Location:            rpath,
				Module:              m,
				Files:               make([]*File, 0),
				ReverseDependencies: make([]*Package, 0),
			}
			m.AddPackage(p)
		}
		// Ensure we know about the astFile in the package.
		// Every astFile is new, so we don't need to check for existence.
		f := p.AddFile(path, astFile)
		ast.Inspect(astFile, func(n ast.Node) bool {
			// Find Import Specs
			imp, ok := n.(*ast.ImportSpec)
			if ok {
				f.AddImport(imp.Path.Value)
			}
			return true
		})
		m.NoOfLines += len(f.Lines)
		return nil
	})
	if err != nil {
		return fmt.Errorf("filepath.Walk: %w", err)
	}
	return nil
}
