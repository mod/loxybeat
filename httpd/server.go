package httpd

import (
	"net/http"
	"time"

	"github.com/mod/loxybeat/config"
)

type Server struct {
	server *http.Server
}

func New(config *config.Config) *Server {
	s := &Server{
		server: &http.Server{
			Addr:           ":8080",
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
	return s
}

func (srv *Server) SetHandleFunc(path string,
	handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(path, handler)
}

func (srv *Server) Start() {
	srv.server.ListenAndServe()
}
