package render

import (
	"context"
	"errors"
	"fmt"
	"github.com/perbu/gogrok/analytics"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	Repo      *analytics.Repo
	srv       *http.Server
	templates map[string]*template.Template
}

func New(repo *analytics.Repo) *Server {
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
	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		panic("Failed to parse frontpage template: " + err.Error())
	}
	s.templates["index"] = tmpl
	return s
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
