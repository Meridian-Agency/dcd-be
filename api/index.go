package handler

import (
	"fmt"
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
	// Temporary debug helper to inspect Vercel request variables
	if r.URL.Query().Get("debug") == "true" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		headers := make(map[string]string)
		for k, v := range r.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}
		_, _ = w.Write([]byte(fmt.Sprintf(`{"url_path": "%s", "request_uri": "%s", "headers": %+v}`, r.URL.Path, r.RequestURI, headers)))
		return
	}

	// Restore the original path so Gin can route correctly
	if matchedPath := r.Header.Get("X-Matched-Path"); matchedPath != "" {
		r.URL.Path = matchedPath
	} else if u, err := url.Parse(r.RequestURI); err == nil {
		r.URL.Path = u.Path
	}

	app.ServeHTTP(w, r)
}
