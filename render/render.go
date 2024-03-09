package render

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
)

type Frontpage struct {
	LocalModules    []Module
	ExternalModules []Module
}

type Module struct {
	Path                  string
	Dependencies          []*Module
	NoOfDependencies      int
	ReverseDependencies   []*Module
	NoOfReverseDependents int
	Packages              []*Package
}

type Package struct {
	Path                  string
	Dependencies          []*Package
	NoOfDependencies      int
	ReverseDependencies   []*Package
	NoOfReverseDependents int
	Files                 []*File
	NofFiles              int
}

type File struct {
	Path     string
	Package  *Package
	Module   *Module
	Lines    []string
	NofLines int
}

//go:embed frontpage.gohtml
var frontpageTemplate string

var templates map[string]*template.Template

func init() {
	templates = make(map[string]*template.Template)
	tmpl, err := template.New("mainPage").Parse(frontpageTemplate)
	if err != nil {
		panic("Failed to parse frontpage template: " + err.Error())
	}
	templates["frontpage"] = tmpl
}

func RenderFrontpage(fp Frontpage, target string) error {
	// open the target file:
	file, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	defer file.Close()
	// render fp
	err = templates["frontpage"].Execute(file, fp)
	if err != nil {
		return fmt.Errorf("templates[\"frontpage\"].Execute: %w", err)
	}
	return nil
}
