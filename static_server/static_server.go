package staticserver

import (
	"net/http"

	"github.com/osamikoyo/adori/config"
	"github.com/osamikoyo/adori/core"
	"github.com/osamikoyo/adori/logger"
)

type StaticServer struct{
	httpserver *http.Server
	logger *logger.Logger
}

func NewStaticServer(
	cfg *config.Config,
	logger *logger.Logger,
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
		logger: logger,
		httpserver: server,
	}
}