package httpd

import (
	"html"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/elastic/beats/libbeat/logp"

	"github.com/mod/loxybeat/config"
)

type Pipeline chan []byte

type Server struct {
	listener net.Listener
	server   *http.Server
	config   *config.Config
	Pipe     Pipeline
	stop     chan bool
}

type LogMux struct {
	logname string
	pipe    Pipeline
}

func (lx *LogMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := html.EscapeString(r.URL.Path)
	logp.Debug("httpd", "Request %s %s", r.Method, path)

	defer r.Body.Close()
	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logp.Err("Read error:", err)
	}
	select {
	case lx.pipe <- buffer:
		logp.Debug("httpd", "Payload sent: %s", buffer)
	default:
		logp.Warn("Pipeline is full discarding log")
	}
	w.Write([]byte("Payload processed\n"))
}

func New(config *config.Config) *Server {
	pipe := make(Pipeline, config.QueueSize)
	http := &http.Server{
		Addr:           config.Address,
		Handler:        &LogMux{pipe: pipe},
		ReadTimeout:    config.Timeout,
		WriteTimeout:   config.Timeout,
		MaxHeaderBytes: 1 << 20,
	}
	srv := &Server{
		server: http,
		config: config,
		Pipe:   pipe,
		stop:   make(chan bool),
	}
	logp.Info("Server Address (%s) - Queue size (%d) - Timeout (%d)",
		config.Address, config.QueueSize, config.Timeout)
	return srv
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
