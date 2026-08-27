package handler

import (
	"net/http"
	"net/url"

	"dcd-be/internal/config"
	"dcd-be/internal/server"
)

var app http.Handler

func init() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	app = server.NewHandler(cfg)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Restore the original path so Gin can route correctly
	if matchedPath := r.Header.Get("X-Matched-Path"); matchedPath != "" {
		r.URL.Path = matchedPath
	} else if u, err := url.Parse(r.RequestURI); err == nil {
		r.URL.Path = u.Path
	}

	app.ServeHTTP(w, r)
}
