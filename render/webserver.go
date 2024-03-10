package render

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/perbu/gogrok/analytics"
	"html/template"
	"io/fs"
	"log/slog"
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
	templates *template.Template
}

func New(repo *analytics.Repo) (*Server, error) {
	const defaultPort = 8080
	s := &Server{
		Repo: repo,
	}
	r := loggingMiddleware(makeMux(s))
	addr := fmt.Sprintf(":%d", getEnvInt("PORT", defaultPort))
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	s.srv = srv

	tmplts, err := loadTemplates()
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
func loadTemplates() (*template.Template, error) {
	tmpl := template.New("")
	err := fs.WalkDir(assets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".gohtml") {
			fileData, err := fs.ReadFile(assets, path)
			if err != nil {
				return err
			}
			_, err = tmpl.New(path).Parse(string(fileData))
			if err != nil {
				return err
			}
			slog.Info("loaded template", "path", path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("fs.WalkDir: %w", err)
	}
	return tmpl, nil
}
