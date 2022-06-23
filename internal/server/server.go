package server

import (
	"net/http"

	"github.com/giffone/forum-security/internal/config"
)

func NewServerTLS(mux *http.ServeMux, port string) *http.Server {
	return &http.Server{
		Addr:         port,
		Handler:      mux,
		TLSConfig:    config.NewCrtConf(),
		ReadTimeout:  config.TimeLimit5s,
		WriteTimeout: config.TimeLimit10s,
		IdleTimeout:  config.TimeLimit20s,
	}
}

func NewServer(mux *http.ServeMux, port string) *http.Server {
	return &http.Server{
		Addr:    port,
		Handler: mux,
	}
}
