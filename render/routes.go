package render

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"time"
)

func makeMux(s *Server) *mux.Router {
	gmux := mux.NewRouter()
	gmux.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	// make a 404 handler:
	// mux.HandleFunc("/", notFoundHandler)
	gmux.HandleFunc("/", makeStaticHandler("assets/index.html")).Methods(http.MethodGet)
	// serve the styles.css directly from the assets embedded filesystem:
	gmux.HandleFunc("/styles.css", makeStaticHandler("assets/styles.css")).Methods(http.MethodGet)
	gmux.HandleFunc("/local", s.handleLocalModuleList).Methods(http.MethodGet)
	gmux.HandleFunc("/external", s.handleExternalModuleList).Methods(http.MethodGet)
	gmux.HandleFunc("/about", s.handleAbout).Methods(http.MethodGet)
	gmux.HandleFunc("/module/{module:.*}", s.handleModule).Methods(http.MethodGet)
	gmux.HandleFunc("/package/{module:[^?]*}", s.handlePackage).Methods(http.MethodGet)
	gmux.HandleFunc("/file/{module:[^?]*}", s.handleFile).Methods(http.MethodGet)
	return gmux
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &responseObserver{ResponseWriter: w}
		next.ServeHTTP(lrw, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "status", lrw.status,
			"duration", time.Since(start), "written", lrw.written)
	})
}

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}
