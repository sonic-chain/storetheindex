package httpfinderserver

import (
	"context"
	"fmt"
	"net"
	"net/http"

	indexer "github.com/filecoin-project/go-indexer-core"
	"github.com/filecoin-project/storetheindex/internal/providers"
	"github.com/filecoin-project/storetheindex/server/finder/http/handler"
	"github.com/gorilla/mux"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("findserver")

type Server struct {
	server *http.Server
	l      net.Listener
}

func (s *Server) URL() string {
	return fmt.Sprint("http://", s.l.Addr().String())
}

func New(listen string, indexer indexer.Interface, registry *providers.Registry, options ...ServerOption) (*Server, error) {
	var cfg serverConfig
	if err := cfg.apply(options...); err != nil {
		return nil, err
	}
	var err error

	l, err := net.Listen("tcp", listen)
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter().StrictSlash(true)
	server := &http.Server{
		Handler:      r,
		WriteTimeout: cfg.apiWriteTimeout,
		ReadTimeout:  cfg.apiReadTimeout,
	}
	s := &Server{server, l}

	// Resource handlers
	h := handler.New(indexer, registry)

	// Client routes
	r.HandleFunc("/cid/{cid}", h.GetSingleCid).Methods("GET")
	r.HandleFunc("/cid", h.GetBatchCid).Methods("POST")

	return s, nil
}

func (s *Server) Start() error {
	log.Infow("finder listening", "addr", s.l.Addr())
	return s.server.Serve(s.l)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}