package render

import (
	"github.com/gorilla/mux"
	"github.com/perbu/gogrok/analytics"
	"github.com/perbu/gogrok/render/fragments"
	"log/slog"
	"net/http"
	"strings"
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

func (s *Server) handleLocalModuleList(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("search")
	mods := s.Repo.ModuleFilter(analytics.DepTypeLocal, filter)
	filteredMods := make([]analytics.Module, 0)
	for _, mod := range mods {
		if strings.Contains(mod.Path, filter) {
			filteredMods = append(filteredMods, mod)
		}
	}
	err := fragments.LocalModules(filteredMods).Render(r.Context(), w)
	if err != nil {
		slog.Error("templ Render", "fragment", "localModules",
			"filter", filter, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) handleExternalModuleList(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("search")
	mods := s.Repo.ModuleFilter(analytics.DepTypeExternal, filter)
	filteredMods := make([]analytics.Module, 0)
	for _, mod := range mods {
		if strings.Contains(mod.Path, filter) {
			filteredMods = append(filteredMods, mod)
		}
	}
	err := fragments.ExternalModules(filteredMods).Render(r.Context(), w)
	if err != nil {
		slog.Error("templ Render", "fragment", "externalModules",
			"filter", filter, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) handleAbout(w http.ResponseWriter, r *http.Request) {
	// err := s.templates.ExecuteTemplate(w, "about", nil)
	err := fragments.About().Render(r.Context(), w)
	if err != nil {
		slog.Error("templ Render", "fragment", "about", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) handleModule(writer http.ResponseWriter, request *http.Request) {
	// get the module name from the URL
	vars := mux.Vars(request)
	moduleName := vars["module"]
	slog.Info("handleModule", "module", moduleName)
	mod, ok := s.Repo.GetModule(moduleName)
	if !ok {
		http.Error(writer, "module not found", http.StatusNotFound)
		return
	}
	err := fragments.Module(mod).Render(request.Context(), writer)
	if err != nil {
		slog.Error("templ Render", "fragment", "module", "module", moduleName, "error", err)
		http.Error(writer, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) handlePackage(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	moduleName := vars["module"]
	packageName := request.URL.Query().Get("package")
	slog.Info("handlePackage", "module", moduleName, "package", packageName)
	mod, ok := s.Repo.GetModule(moduleName)
	if !ok {
		http.Error(writer, "module not found", http.StatusNotFound)
		return
	}
	// find the package in the module:
	var pkg *analytics.Package
	for _, p := range mod.Packages {
		if p.Name == packageName {
			pkg = p
			break
		}
	}
	if pkg == nil {
		http.Error(writer, "package not found", http.StatusNotFound)
		return
	}

	err := fragments.Package(pkg).Render(request.Context(), writer)
	if err != nil {
		slog.Error("templ Render", "fragment", "package", "module", moduleName, "package", packageName, "error", err)
		http.Error(writer, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) handleFile(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	moduleName := vars["module"]
	packageName := request.URL.Query().Get("package")
	fileName := request.URL.Query().Get("file")
	slog.Info("handleFile", "module", moduleName, "package", packageName, "file", fileName)
	// find the module in the repo:
	mod, ok := s.Repo.GetModule(moduleName)
	if !ok {
		http.Error(writer, "module not found", http.StatusNotFound)
		return
	}
	// find the package and the file:
	var file *analytics.File
	for _, pkg := range mod.Packages {
		if pkg.Name == packageName {
			for _, f := range pkg.Files {
				if f.Name == fileName {
					file = f
					break
				}
			}
		}
	}
	if file == nil {
		http.Error(writer, "file not found", http.StatusNotFound)
		return
	}
	err := fragments.File(file).Render(request.Context(), writer)
	if err != nil {
		slog.Error("templ Render", "fragment", "file", "module", moduleName, "package", packageName, "file", fileName, "error", err)
		http.Error(writer, "internal server error", http.StatusInternalServerError)
	}
}
