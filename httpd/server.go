package httpd

import (
	"net"
	"net/http"

	"github.com/elastic/beats/libbeat/logp"

	"github.com/mod/loxybeat/config"
)

type Server struct {
	listener net.Listener
	server   *http.Server
	config   *config.Config
	stop     chan bool
}

func New(config *config.Config) *Server {
	http := &http.Server{
		Addr:           config.Address,
		ReadTimeout:    config.Timeout,
		WriteTimeout:   config.Timeout,
		MaxHeaderBytes: 1 << 20,
	}
	srv := &Server{
		server: http,
		config: config,
		stop:   make(chan bool),
	}
	return srv
}

func (srv *Server) SetHandleFunc(path string,
	handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(path, handler)
}

func (srv *Server) Start() {
	var err error

	logp.Info("HTTP listening on TCP %s", srv.config.Address)

	srv.listener, err = net.Listen("tcp", srv.config.Address)
	if err != nil {
		logp.Err("HTTP failed to start TCP4 listener: %v", err)
		return
	}

	go srv.serve()

	// Wait for shutdown
	select {
	case <-srv.stop:
		logp.Info("HTTP server shutting down on request")
		srv.close()
	}
}

func (srv *Server) Stop() {
	logp.Info("httpd.Stop() Recieved")
	srv.stop <- true
}

func (srv *Server) serve() {
	err := srv.server.Serve(srv.listener)

	select {
	case <-srv.stop:
		logp.Info("Server shutdown")

	default:
		logp.Err("HTTP server failed: %v", err)
		srv.close()
		return
	}
}

func (srv *Server) close() {
	if err := srv.listener.Close(); err != nil {
		logp.Err("Failed to close HTTP listener: %v", err)
	}
}
