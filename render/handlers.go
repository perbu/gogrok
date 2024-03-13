package render

import (
	"github.com/gorilla/mux"
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

func (s *Server) handleLocalModuleList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")
	overview := make([]ModuleOverviewModule, 0)
	for _, mod := range s.Repo.ModuleFilter(analytics.DepTypeLocal, query) {
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
		deps := make([]string, 0)
		for _, dep := range mod.Dependencies {
			deps = append(deps, dep.Path)
		}
		revDeps := make([]string, 0)
		for _, dep := range mod.ReverseModuleDependencies {
			revDeps = append(revDeps, dep.Path)
		}
		overview = append(overview, ModuleOverviewModule{
			Path:          mod.Path,
			Version:       mod.Version,
			PackagesCount: len(mod.Packages),
			FilesCount:    files,
			LinesOfCode:   lines,
			Packages:      packages,
			Deps:          deps,
			RevDeps:       revDeps,
		})
	}
	mo := ModuleOverview{Modules: overview}
	err := s.templates.ExecuteTemplate(w, "localModules", mo)
	if err != nil {
		slog.Error("execute", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) handleExternalModuleList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")
	overview := make([]ModuleOverviewModule, 0)
	for _, mod := range s.Repo.ModuleFilter(analytics.DepTypeExternal, query) {
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
		revDeps := make([]string, 0)
		for _, dep := range mod.ReverseModuleDependencies {
			revDeps = append(revDeps, dep.Path)
		}
		overview = append(overview, ModuleOverviewModule{
			Path:          mod.Path,
			Version:       mod.Version,
			PackagesCount: len(mod.Packages),
			FilesCount:    files,
			LinesOfCode:   lines,
			Packages:      packages,
			RevDeps:       revDeps,
		})
	}
	mo := ModuleOverview{Modules: overview}
	err := s.templates.ExecuteTemplate(w, "externalModules", mo)
	if err != nil {
		slog.Error("execute", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) handleAbout(w http.ResponseWriter, r *http.Request) {
	err := s.templates.ExecuteTemplate(w, "about", nil)
	if err != nil {
		slog.Error("execute", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) handleModule(writer http.ResponseWriter, request *http.Request) {
	// get the module name from the URL
	vars := mux.Vars(request)
	moduleName := vars["path"]
	slog.Info("handleModule", "module", moduleName)
	mod, ok := s.Repo.GetModule(moduleName)
	if !ok {
		http.Error(writer, "module not found", http.StatusNotFound)
		return
	}
	// render the module template. First make a ModuleDetailModule from the module
	// and then execute the template with it.
	mdm := mod2mod(mod)
	err := s.templates.ExecuteTemplate(writer, "module", mdm)
	if err != nil {
		slog.Error("execute", "error", err)
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
	// render the module template. First make a ModuleDetailModule from the module
	// and then execute the template with it.
	p := package2package(mod, pkg)
	err := s.templates.ExecuteTemplate(writer, "package", p)
	if err != nil {
		slog.Error("execute", "error", err)
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
	// find the file in the module:
	var file *analytics.File
	for _, pkg := range mod.Packages {
		for _, f := range pkg.Files {
			if f.Name == fileName {
				file = f
				break
			}
		}
	}
	if file == nil {
		http.Error(writer, "file not found", http.StatusNotFound)
		return
	}
	// render the file template. First make a File from the file
	// and then execute the template with it.
	f := file2file(mod, file)
	err := s.templates.ExecuteTemplate(writer, "file", f)
	if err != nil {
		slog.Error("execute", "error", err)
		http.Error(writer, "internal server error", http.StatusInternalServerError)
	}
}

func package2package(mod *analytics.Module, pkg *analytics.Package) Package {
	files := make([]PackageFile, 0)
	for _, file := range pkg.Files {
		pf := PackageFile{
			Name:    file.Name,
			Package: pkg.Name,
			Module:  mod.Path,
		}
		files = append(files, pf)
	}
	return Package{
		Name:     pkg.Name,
		Location: pkg.Location,
		Module:   mod.Path,
		Files:    files,
	}
}

func mod2mod(mod *analytics.Module) ModuleDetailModule {
	noOfLines := 0
	noOfFiles := 0
	packages := make([]ModuleDetailModulePackage, 0)
	for _, pkg := range mod.Packages {
		mdmp := ModuleDetailModulePackage{
			Name:   pkg.Name,
			Module: mod.Path,
		}
		packages = append(packages, mdmp)
		for _, file := range pkg.Files {
			noOfLines += len(file.Lines)
			noOfFiles++
		}
	}
	return ModuleDetailModule{
		Path:          mod.Path,
		Version:       mod.Version,
		Packages:      packages,
		PackagesCount: len(mod.Packages),
		FilesCount:    noOfFiles,
		LinesOfCode:   noOfLines,
	}
}
