package gateway

import (
	"context"
	"net/http"

	"github.com/osamikoyo/adori/config"
	"github.com/osamikoyo/adori/core"
	"github.com/osamikoyo/adori/logger"
	"github.com/osamikoyo/adori/proxy"
	"go.uber.org/zap"
)

type GatewayServer struct {
	httpserver *http.Server
	logger     *logger.Logger
}

func NewGatewayServer(
	cfg *config.Config,
	logger *logger.Logger,
	core *core.AdoriCore,
) (*GatewayServer, error) {
	proxy, err := proxy.NewProxy(cfg)
	if err != nil{
		logger.Fatal("failed get proxy server", zap.Error(err))

		return nil, err
	}

	return &GatewayServer{
		httpserver: &http.Server{
			Addr: cfg.Addr,
			Handler: core.CoreMiddlewareForHandlerFunc(proxy.ServeHTTP),
		},
		logger: logger,
	}, nil
}

func (g *GatewayServer) Run(ctx context.Context) error {
	return g.httpserver.ListenAndServe()
}

func (g *GatewayServer) Stop() error {
	return g.httpserver.Close()
}
