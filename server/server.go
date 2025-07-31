package server

import (
	"github.com/osamikoyo/adori/gateway"
	"github.com/osamikoyo/adori/logger"
	staticserver "github.com/osamikoyo/adori/static_server"
)

type Server struct {
	staticserver *staticserver.StaticServer
	gateway *gateway.GatewayServer
	logger *logger.Logger
}
