package render

import (
	"github.com/perbu/gogrok/analytics"
	"log/slog"
	"net/http"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	slog.Warn("not found", "path", r.URL.Path)
	http.Error(w, "not found", http.StatusNotFound)
}

func makeStaticHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("serving static file", "path", path)
		data, err := assets.ReadFile(path)
		if err != nil {
			slog.Error("assets.ReadFile", "path", path, "error", err)
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		_, _ = w.Write(data)
	}
}

func (s *Server) handleLocalModuleOverview(w http.ResponseWriter, r *http.Request) {
	overview := make([]ModuleOverviewModule, 0)
	for _, mod := range s.Repo.ModuleFilter(analytics.DepTypeLocal) {
		lines := 0
		files := 0
		packages := make([]string, 0)
		for _, pkg := range mod.Packages {
			packages = append(packages, pkg.Name)
			for _, file := range pkg.Files {
				lines += len(file.Lines)
				files++
			}
		}
		overview = append(overview, ModuleOverviewModule{
			Path:          mod.Path,
			Version:       mod.Version,
			PackagesCount: len(mod.Packages),
			FilesCount:    files,
			LinesOfCode:   lines,
			Packages:      packages,
		})
	}
	mo := ModuleOverview{Modules: overview}
	err := s.templates.ExecuteTemplate(w, "localModules", mo)
	if err != nil {
		slog.Error("execute", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

}
