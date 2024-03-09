package render

import (
	"cmp"
	_ "embed"
	"github.com/perbu/gogrok/analytics"
	"net/http"
	"slices"
)

//go:embed index.gohtml
var indexTemplate string

// handleFrontpage is the handler for the frontpage of the webserver, uses the "index"
// template to render the frontpage
func (s *Server) handleFrontpage(w http.ResponseWriter, r *http.Request) {
	lm := s.Repo.ModuleFilter(analytics.DepTypeLocal)
	slices.SortFunc(lm, func(a, b analytics.Module) int {
		return cmp.Compare(a.Path, b.Path)
	})
	em := s.Repo.ModuleFilter(analytics.DepTypeExternal)
	slices.SortFunc(em, func(a, b analytics.Module) int {
		return cmp.Compare(a.Path, b.Path)
	})
	rlm := make([]Module, 0, len(lm))
	for _, m := range lm {
		rlm = append(rlm, Module{
			Path:                m.Path,
			Dependencies:        m.GetStringDependencies(),
			ReverseDependencies: m.GetStringReverseDependencies(),
		})
	}
	elm := make([]Module, 0, len(em))
	for _, m := range em {
		elm = append(elm, Module{
			Path:                m.Path,
			Dependencies:        m.GetStringDependencies(),
			ReverseDependencies: m.GetStringReverseDependencies(),
		})
	}
	err := s.templates["index"].Execute(w, Frontpage{
		LocalModules:    rlm,
		ExternalModules: elm,
	})
	if err != nil {
		http.Error(w, "Failed to render frontpage: "+err.Error(), http.StatusInternalServerError)
	}
}

//go:embed style.css
var styleBytes []byte

// styleHandler is the handler for the /style.css file, serves the style.css file
func styleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	_, _ = w.Write(styleBytes)
}
