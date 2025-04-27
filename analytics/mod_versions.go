package analytics

import "golang.org/x/mod/semver"

func (m *Module) IsLatestVersion(version string) bool {
	semver.Sort(m.versions)
	latest := m.versions[len(m.versions)-1]
	return version == latest
}
