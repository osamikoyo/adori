package staticserver

import (
	"context"
	"net/http"

	"github.com/osamikoyo/adori/config"
	"github.com/osamikoyo/adori/core"
)

type StaticServer struct{
	httpserver *http.Server
}

func NewStaticServer(
	cfg *config.Config,
	core *core.AdoriCore,
) *StaticServer {
	fs := http.FileServer(http.Dir(cfg.Static.Dir))
	
	mux := http.NewServeMux()

	mux.Handle(cfg.Static.Prefix, core.CoreMiddlewareForHandler(http.StripPrefix(cfg.Static.Prefix, fs),))

	server := &http.Server{
		Addr: cfg.Addr,
		Handler: mux,
	}

	return &StaticServer{
		httpserver: server,
	}
}

func (s *StaticServer) Run(ctx context.Context) error {
	return s.httpserver.ListenAndServe()
}

func (s *StaticServer) Stop() error {
	return s.httpserver.Close()
}