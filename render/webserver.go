package render

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/perbu/gogrok/analytics"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//go:embed assets/*
var assets embed.FS

type Server struct {
	Repo      *analytics.Repo
	srv       *http.Server
	templates map[string]*template.Template
}

func New(repo *analytics.Repo) (*Server, error) {
	const defaultPort = 8080
	s := &Server{
		Repo:      repo,
		templates: make(map[string]*template.Template),
	}
	r := makeMux(s)
	addr := fmt.Sprintf(":%d", getEnvInt("PORT", defaultPort))
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	s.srv = srv

	tmplts, err := loadTemplates("assets")
	if err != nil {
		panic("Failed to load templates: " + err.Error())
	}
	s.templates = tmplts
	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		timeout, cancel := context.WithTimeout(context.Background(), 5)
		defer cancel()
		_ = s.srv.Shutdown(timeout)
	}()
	err := s.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("srv.ListenAndServe: %w", err)
	}
	return nil
}

func getEnvInt(key string, defaultValue int) int {
	str, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return val
}

// loadTemplates loads templates from the embedded file system into a map.
func loadTemplates(dir string) (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)

	// Read the directory from the embedded file system.
	templateFiles, err := fs.ReadDir(assets, dir)
	if err != nil {
		return nil, fmt.Errorf("fs.ReadDir: %w", err)
	}

	// Range over the files and parse them as templates.
	for _, entry := range templateFiles {
		// Skip directories and non-gohtml files.
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".gohtml") {
			continue
		}
		fileName := entry.Name()
		path := dir + "/" + fileName
		tmpl, err := template.ParseFS(assets, path)
		if err != nil {
			return nil, fmt.Errorf("template.ParseFS: %w", err)
		}
		templates[fileName] = tmpl
	}
	return templates, nil
}
