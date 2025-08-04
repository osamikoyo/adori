package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/osamikoyo/adori/cash"
	"github.com/osamikoyo/adori/config"
	"github.com/osamikoyo/adori/core"
	"github.com/osamikoyo/adori/defence"
	"github.com/osamikoyo/adori/gateway"
	"github.com/osamikoyo/adori/logger"
	"github.com/osamikoyo/adori/models"
	staticserver "github.com/osamikoyo/adori/static_server"
	"github.com/osamikoyo/adori/statistic"
	"go.uber.org/zap"
)

type (
	Subserver interface {
		Run(ctx context.Context) error
		Stop() error
	}

	Server struct {
		subserver Subserver
		logger    *logger.Logger
	}
)

func NewAdoriServer(cfg *config.Config, logger *logger.Logger) (*Server, error) {
	logger.Info("starting adori...")

	var statisticChan chan *models.StatisticChunk

	switch cfg.Regime {
	case "gateway":
		logger.Info("prepare gateway server", zap.String("addr", cfg.Addr))

		core := core.NewAdoriCore(
			cash.NewLocalCash(time.Duration(cfg.Cash.IntervalInSeconds)*time.Second, logger),
			defence.NewDefence(cfg, logger),
			statistic.NewStatisticClient(statisticChan),
			logger,
		)

		gateway, err := gateway.NewGatewayServer(cfg, logger, core)
		if err != nil {
			logger.Fatal("failed get gateway server", zap.Error(err))

			return nil, err
		}

		return &Server{
			subserver: gateway,
			logger:    logger,
		}, nil
	case "static":
		logger.Info("prepare static server", zap.String("addr", cfg.Addr))

		core := core.NewAdoriCore(
			cash.NewLocalCash(time.Duration(cfg.Cash.IntervalInSeconds)*time.Second, logger),
			defence.NewDefence(cfg, logger),
			statistic.NewStatisticClient(statisticChan),
			logger,
		)
		
		static := staticserver.NewStaticServer(cfg, core)

		return &Server{
			subserver: static,
			logger: logger,
		}, nil
		default:
			logger.Fatal("unknown regime", zap.String("regime", cfg.Regime))

			return nil, fmt.Errorf("unknown regime: %s", cfg.Regime)
	}
}

func (s *Server) Run(ctx context.Context) error {
	err := s.subserver.Run(ctx)
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	return s.subserver.Stop()
}