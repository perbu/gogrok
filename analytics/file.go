package analytics

import (
	"bufio"
	"go/ast"
	"go/token"
	"log/slog"
	"os"
	"path"
	"strings"
)

// AddFile adds a file to the package
func (p *Package) AddFile(name string, file *ast.File) *File {
	for _, file := range p.files {
		if file.Name == name {
			// should never happen
			panic("file already exists")
		}
	}
	//
	lines, err := readFile(name)
	if err != nil {
		// should never happen
		panic("readFile: " + err.Error())

	}
	fileType := NameToFileType(name)
	// try to detect if the file is generated

	if fileType == HumanGo && len(lines) > 5 {
		for _, line := range lines[:5] {
			if strings.Contains(line, "Code generated") {
				fileType = GeneratedGo
				break
			}
		}
	}

	// create a new file:
	f := &File{
		Name:    path.Base(name),
		Imports: make([]*Package, 0),
		Package: p,
		Module:  p.Module,
		source:  lines,
		ast:     file,
		Type:    fileType,
	}
	p.files = append(p.files, f)
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
	// ignore if this in an internal import within the same module:
	if strings.HasPrefix(name, f.Module.Path) {
		return
	}
	// find the module name:
	mName, ok := f.Module.Repo.FindModule(name)
	if !ok {
		// The modules we can't locate are either stdlib or external, not our code.
		return
	}
	mod, ok := f.Module.Repo.GetModule(mName)
	if !ok {
		panic("GetModule failed " + mName)
	}
	// skip if the package isn't local:
	if mod.Type != DepTypeLocal {
		return
	}
	pack := packageNameFromImportPath(name)
	p, ok := mod.GetPackage(pack)
	if !ok {
		// package is liked missing from the latest versions of the module, we just ignore it
		return
	}
	f.Imports = append(f.Imports, p)
}

func packageNameFromImportPath(path string) string {
	// package name is the last part of the import path
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

// readFile reads a file and return a list of lines
func readFile(path string) ([]string, error) {
	// Open the file for reading.
	file, err := os.Open(path)
	if err != nil {
		return nil, err // Return the error if the file cannot be opened.
	}
	defer file.Close() // Ensure the file is closed after this function completes.

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()) // Add the line to the lines slice.
	}

	if err := scanner.Err(); err != nil {
		return nil, err // Return the error if a scanning error occurs.
	}

	return lines, nil // Return the slice of lines and a nil error.
}

func (f *File) Lines() int {
	return len(f.source)
}

func (f *File) GetSource() []string {
	return f.source
}

// complexityVisitor implements the ast.Visitor interface, counting decision points.
type complexityVisitor struct {
	complexity int
}

// Visit inspects the AST nodes and calculates complexity based on certain node types.
func (v *complexityVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch n.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
		v.complexity++
	case *ast.CaseClause:
		// Each case in a switch is a new path
		v.complexity++
	case *ast.BinaryExpr:
		// Increase complexity for logical AND and OR operations
		be := n.(*ast.BinaryExpr)
		if be.Op == token.LAND || be.Op == token.LOR {
			v.complexity++
		}
	}

	return v
}

// CalculateComplexity walks the AST of a Go file to calculate its cyclomatic complexity.
func (f *File) CalculateComplexity() float32 {
	if f.Type == GeneratedGo {
		slog.Info("complex skipping generated file", "file", f.Name)
		return 1.0
	}
	v := &complexityVisitor{}
	ast.Walk(v, f.ast)
	return float32(v.complexity) + 1 // Adding 1 for the entry point
}

func NameToFileType(name string) FileType {
	switch path.Ext(name) {
	case ".go":
		// check if the name ends in .gen.go
		if strings.HasSuffix(name, ".gen.go") {
			return GeneratedGo
		}
		if strings.HasSuffix(name, "_test.go") {
			return TestGo
		}
		return HumanGo
	default:
		return Unknown
	}
}
