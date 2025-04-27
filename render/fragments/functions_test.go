package fragments

import (
	"github.com/perbu/gogrok/analytics"
	"testing"
)

func TestS(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "int",
			input:    42,
			expected: "42",
		},
		{
			name:     "int64",
			input:    int64(42),
			expected: "42",
		},
		{
			name:     "float",
			input:    42.5,
			expected: "42.5",
		},
		{
			name:     "string",
			input:    "hello",
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s(tt.input)
			if result != tt.expected {
				t.Errorf("s(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSlen(t *testing.T) {
	tests := []struct {
		name     string
		input    []any
		expected string
	}{
		{
			name:     "empty slice",
			input:    []any{},
			expected: "0",
		},
		{
			name:     "slice with one element",
			input:    []any{1},
			expected: "1",
		},
		{
			name:     "slice with multiple elements",
			input:    []any{1, 2, 3},
			expected: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slen(tt.input)
			if result != tt.expected {
				t.Errorf("slen(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestModuleUrl(t *testing.T) {
	tests := []struct {
		name     string
		module   *analytics.Module
		expected string
	}{
		{
			name: "simple module",
			module: &analytics.Module{
				Path: "github.com/example/module",
			},
			expected: "/module/github.com/example/module",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := moduleUrl(tt.module)
			if result != tt.expected {
				t.Errorf("moduleUrl(%v) = %v, expected %v", tt.module.Path, result, tt.expected)
			}
		})
	}
}

func TestPackageUrl(t *testing.T) {
	tests := []struct {
		name     string
		pkg      *analytics.Package
		expected string
	}{
		{
			name: "simple package",
			pkg: &analytics.Package{
				Name: "pkg",
				Module: &analytics.Module{
					Path: "github.com/example/module",
				},
			},
			expected: "/package/github.com/example/module?package=pkg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := packageUrl(tt.pkg)
			if result != tt.expected {
				t.Errorf("packageUrl(%v) = %v, expected %v", tt.pkg.Name, result, tt.expected)
			}
		})
	}
}

func TestFileUrl(t *testing.T) {
	tests := []struct {
		name     string
		file     *analytics.File
		expected string
	}{
		{
			name: "simple file",
			file: &analytics.File{
				Name: "file.go",
				Package: &analytics.Package{
					Name: "pkg",
				},
				Module: &analytics.Module{
					Path: "github.com/example/module",
				},
			},
			expected: "/file/github.com/example/module?file=file.go&package=pkg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fileUrl(tt.file)
			if result != tt.expected {
				t.Errorf("fileUrl(%v) = %v, expected %v", tt.file.Name, result, tt.expected)
			}
		})
	}
}
