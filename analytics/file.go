package analytics

import (
	"bufio"
	"go/ast"
	"os"
	"path"
	"strings"
)

// AddFile adds a file to the package
func (p *Package) AddFile(name string, file *ast.File) *File {
	for _, file := range p.Files {
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

	// create a new file:
	f := &File{
		Name:    path.Base(name),
		Imports: make([]*Package, 0),
		Package: p,
		Module:  p.Module,
		Lines:   lines,
		Ast:     file,
	}
	p.Files = append(p.Files, f)
	p.NoOfLines += len(lines)
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
