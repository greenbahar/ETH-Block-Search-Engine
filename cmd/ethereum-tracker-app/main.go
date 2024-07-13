package main

import (
	"context"
	"ethereum-tracker-app/cmd/config"
	routers "ethereum-tracker-app/internal/http"
	"ethereum-tracker-app/internal/http/handlers"
	"ethereum-tracker-app/internal/services/blockprocessor"
	"ethereum-tracker-app/internal/services/blocksearch"
	"ethereum-tracker-app/internal/storage/inmemorydb"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	ctx := context.Background()
	logger := log.New(os.Stdout, "app: ", log.LstdFlags)
	if err := godotenv.Load(); err != nil {
		logger.Fatal("No .env file found")
	}
	systemConfig := config.LoadConfig(logger)

	//todo: variadic function to pass ...options to constructors. passing option functions instead of one-by-one entities
	storage := inmemorydb.NewInmemortDBService(*systemConfig, logger)
	blockprocessService := blocksearch.NewServie(*systemConfig, logger, storage)
	handler := handlers.NewHandler(blockprocessService)
	router := routers.SetupRouters(handler)
	ethClient, ethClientErr := blockprocessor.NewEthClient(ctx, *systemConfig, logger, storage)
	if ethClientErr != nil {
		logger.Fatal(errors.Wrap(ethClientErr, "cannot stablish Ethereum client"))
	}

	appService := &Service{
		Config:            systemConfig,
		Logger:            logger,
		EthClient:         ethClient,
		BlockProcessSvc:   blockprocessService,
		Router:            router,
		InMemoryDBService: storage,
		Server: &http.Server{
			Handler:      router,
			Addr:         fmt.Sprintf("%s:%s", systemConfig.ServerConf.ServerIP, systemConfig.ServerConf.ServerPort),
			WriteTimeout: systemConfig.ServerConf.WriteTimeout,
			ReadTimeout:  systemConfig.ServerConf.ReadTimeout,
		},
	}

	if err := appService.run(ctx); err != nil {
		appService.Logger.Fatalf("Failed to run the service: %v", err)
	}
}
