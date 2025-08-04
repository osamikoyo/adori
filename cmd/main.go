package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/osamikoyo/adori/config"
	"github.com/osamikoyo/adori/logger"
	"github.com/osamikoyo/adori/server"
	"go.uber.org/zap"
)


var ConfigfileNames = []string{"adori.yml", "adori.yaml", "adori-config.yaml", "adori-config.yml"}

func main() {
	art := `
   ___     __         _ 
  / _ |___/ /__  ____(_)
 / __ / _  / _ \/ __/ / 
/_/ |_\_,_/\___/_/ /_/  
                        
	`
	fmt.Println(art)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var (
		cfg *config.Config
		err error
	)

	for i, arg := range os.Args{
		if arg == "--config"{
			cfg, err = config.NewConfig(os.Args[i+1])
		}
	}

	if cfg == nil {
		cfg, err = config.NewConfig(ConfigfileNames...)
	}

	if err != nil{
		log.Fatal("failed get config", err)

		return
	}

	logcfg := logger.Config{
		AppName: "adori",
		LogFile: "logs/adori.log",
		LogLevel: "debug",
		AddCaller: false,
	}

	if cfg.Production {
		logcfg.LogLevel = "info"
	}

	if err := logger.Init(logcfg);err != nil{
		log.Fatal(err)
		
		return
	}

	logger := logger.Get()

	logger.Info("preparing adori...")

	logger.Info("config loaded successfully", zap.Any("config", cfg))

	adori, err := server.NewAdoriServer(cfg, logger)
	if err != nil{
		logger.Fatal("failed get adori server", zap.Error(err))

		return
	}

	go func() {
		<- ctx.Done()

		logger.Info("stopping adori...")
		
		adori.Stop()
	}()

	if err = adori.Run(ctx);err != nil{
		logger.Fatal("failed run adori server", zap.Error(err))

		return
	}
}