package render

import (
	_ "embed"
	"fmt"
	"github.com/perbu/gogrok/repo"
	"html/template"
	"os"
)

type Frontpage struct {
	Modules []*repo.Module
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
