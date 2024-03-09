package repo2render

import (
	"github.com/perbu/gogrok/render"
	"github.com/perbu/gogrok/repo"
)

// TransformModule transforms a repo.Module into a render.Module.
func TransformModule(rm *repo.Module) render.Module {
	var rDependencies []render.Module
	for _, dep := range rm.Dependencies {
		rDependencies = append(rDependencies, TransformModule(dep))
	}

	var rReverseDependencies []render.Module
	for _, revDep := range rm.ReverseModuleDependencies {
		rReverseDependencies = append(rReverseDependencies, TransformModule(revDep))
	}

	var rPackages []render.Package
	for _, pkg := range rm.Packages {
		rPackages = append(rPackages, TransformPackage(pkg))
	}

	return render.Module{
		Path:                  rm.Path,
		Dependencies:          rDependencies,
		NoOfDependencies:      len(rDependencies),
		ReverseDependencies:   rReverseDependencies,
		NoOfReverseDependents: len(rReverseDependencies),
		Packages:              rPackages,
	}
}

// TransformPackage transforms a repo.Package into a render.Package.
func TransformPackage(rp *repo.Package) render.Package {
	// Similar transformation logic for packages
	var rDependencies []render.Package
	// Populate dependencies

	var rReverseDependencies []render.Package
	// Populate reverse dependencies

	var rFiles []render.File
	for _, file := range rp.Files {
		rFiles = append(rFiles, TransformFile(file))
	}

	return render.Package{
		Path:                  rp.Location,
		Dependencies:          rDependencies,
		NoOfDependencies:      len(rDependencies),
		ReverseDependencies:   rReverseDependencies,
		NoOfReverseDependents: len(rReverseDependencies),
		Files:                 rFiles,
		NofFiles:              len(rFiles),
	}
}

// TransformFile transforms a repo.File into a render.File.
func TransformFile(rf *repo.File) render.File {
	return render.File{
		Path:     rf.Name,
		Package:  render.Package{}, // Populate if necessary, might require back-reference adjustments
		Module:   render.Module{},  // Populate if necessary, similar concern
		Lines:    rf.Lines,
		NofLines: len(rf.Lines),
	}
}
