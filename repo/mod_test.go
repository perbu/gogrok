package repo

import (
	"testing"
)

func TestFindModule(t *testing.T) {
	// Mock setup
	mockRepo := Repo{
		modules: map[string]*Module{
			"github.com/example/project":     {Path: "github.com/example/project"},
			"github.com/example/project/sub": {Path: "github.com/example/project/sub"},
			"github.com/another/repo":        {Path: "github.com/another/repo"},
		},
	}

	tests := []struct {
		name     string
		path     string
		wantPath string
		wantOk   bool
	}{
		{
			name:     "Direct match",
			path:     "github.com/example/project",
			wantPath: "github.com/example/project",
			wantOk:   true,
		},
		{
			name:     "Match after moving up one level",
			path:     "github.com/example/project/sub/package",
			wantPath: "github.com/example/project/sub",
			wantOk:   true,
		},
		{
			name:     "Match after moving up multiple levels",
			path:     "github.com/example/project/sub/package/module",
			wantPath: "github.com/example/project/sub",
			wantOk:   true,
		},
		{
			name:     "No match found",
			path:     "github.com/unknown/repo",
			wantPath: "",
			wantOk:   false,
		},
		{
			name:     "Highest level without match",
			path:     "github.com",
			wantPath: "",
			wantOk:   false,
		},
		{
			name:     "Standard library import",
			path:     "log/slog",
			wantPath: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotOk := mockRepo.FindModule(tt.path)
			if gotPath != tt.wantPath || gotOk != tt.wantOk {
				t.Errorf("FindModule(%q) = %q, %v; want %q, %v",
					tt.path, gotPath, gotOk, tt.wantPath, tt.wantOk)
			}
		})
	}
}
