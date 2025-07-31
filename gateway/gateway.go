package gateway

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/osamikoyo/adori/logger"
	"go.uber.org/zap"
)

type GatewayServer struct {
	httpserver *http.Server
	proxyTable map[string]*url.URL
	logger     *logger.Logger
}

func NewGatewayServer(
	addr string,
	urltable map[string]string,
	logger *logger.Logger,
) *GatewayServer {
	table := make(map[string]*url.URL)

	mux := http.NewServeMux()

	for prefix, path := range urltable {
		u, err := url.Parse(path)
		if err != nil {
			logger.Warn("failed parse url",
				zap.String("prefix", prefix),
				zap.String("url", path),
				zap.Error(err))
			
			continue
		}

		proxyHandler := httputil.NewSingleHostReverseProxy(u)

		mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
			proxyHandler.ServeHTTP(w, r)
		})
	}

	return &GatewayServer{
		httpserver: &http.Server{
			Addr: addr,
			Handler: mux,
		},
		proxyTable: table,
		logger: logger,
	}
}

func (g *GatewayServer) Run(ctx context.Context) error {
	return g.httpserver.ListenAndServe()
}

func (g *GatewayServer) Close() error {
	return g.httpserver.Close()
}
