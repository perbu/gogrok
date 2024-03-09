package render

import "net/http"

func makeMux(s *Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.handleFrontpage)
	mux.HandleFunc("GET /style.css", styleHandler)
	return mux

}
