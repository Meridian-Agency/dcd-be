package handler

import (
	"net/http"

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
	app.ServeHTTP(w, r)
}
