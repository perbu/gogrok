package analytics

func (r *Repo) reverseDeps() {
	// Initialize a map for reverse package dependencies to avoid duplicates
	packageReverseDepsMap := make(map[*Package]map[*Package]struct{})

	// Iterate through modules in the repo
	for _, module := range r.modules {
		// Populate Reverse Module Dependencies
		for _, dependency := range module.Dependencies {
			dependency.ReverseModuleDependencies = append(dependency.ReverseModuleDependencies, module)
		}

		// Access packages within the module
		for _, pkg := range module.Packages {
			// Ensure the map for the package is initialized
			if _, exists := packageReverseDepsMap[pkg]; !exists {
				packageReverseDepsMap[pkg] = make(map[*Package]struct{})
			}

			// Populate Reverse Package Dependencies using file imports
			for _, file := range pkg.files {
				for _, importedPkg := range file.Imports {
					// Avoid adding duplicate reverse dependencies
					if _, exists := packageReverseDepsMap[importedPkg][pkg]; !exists {
						if _, exists := packageReverseDepsMap[importedPkg]; !exists {
							packageReverseDepsMap[importedPkg] = make(map[*Package]struct{})
						}
						packageReverseDepsMap[importedPkg][pkg] = struct{}{}
						importedPkg.ReverseDependencies = append(importedPkg.ReverseDependencies, pkg)
					}
				}
			}
		}
	}
}
