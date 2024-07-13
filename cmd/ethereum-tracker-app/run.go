package main

import (
	"context"
	"ethereum-tracker-app/cmd/config"
	"ethereum-tracker-app/internal/services/blockprocessor"
	"ethereum-tracker-app/internal/services/blocksearch"
	"ethereum-tracker-app/internal/storage/inmemorydb"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type Service struct {
	Config            *config.Config
	Logger            *log.Logger
	EthClient         blockprocessor.Service
	BlockProcessSvc   blocksearch.Service
	InMemoryDBService inmemorydb.Service
	Router            http.Handler
	Server            *http.Server
}

// run:
//   - starts the service
//   - starts all gouroutines
//   - handles graceful shutdown
func (s *Service) run(ctx context.Context) error {
	blockTxChan := make(chan types.Transactions, s.Config.EthClientConf.NumberOfRecentBlocks)
	wg := &sync.WaitGroup{}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(ctx)

	// Handle OS signals for graceful shutdown
	stopChan := make(chan os.Signal, 1) // a buffer of size 1, so the notifier are not blocked
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stopChan
		s.Logger.Println("Received shutdown signal")
		cancel() // Signal all goroutines to stop
		defer cancel()

		ctxShutdown, cancelServerShutdown := context.WithTimeout(ctx, 5*time.Second)
		defer cancelServerShutdown()
		if err := s.Server.Shutdown(ctxShutdown); err != nil {
			s.Logger.Printf("Error shutting down server: %v", err)
		}
	}()

	// setup a workerpool to process transations in a block concurrently
	for i := 0; i < s.Config.EthClientConf.NumberOfBlockProcessorWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.EthClient.WokerTransactionProcessor(ctx, blockTxChan, wg)
		}()
	}

	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		if err := s.EthClient.SubscribeToNewGeneratedBlocks(ctx); err != nil {
			s.Logger.Fatalf("Failed to subscribe to new generated blocks: %v", err)
		}
	}()

	if err := s.EthClient.FetchAndStoreRecentBlocks(ctx, blockTxChan); err != nil {
		s.Logger.Printf("Failed to fetch and store recent blocks: %v", err)
		return errors.Wrap(err, "failed to fetch and store recent blocks")
	}

	// wait until "FetchAndStoreRecentBlocks" fetches all recent blocks from the blockchain and also finish processing them at "WokerTransactionProcessor" workers
	wg.Wait()

	s.Logger.Printf("Server is running on %s", s.Server.Addr)
	err := s.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.Logger.Printf("error in ListenAndServe")

		return err
	} else if err == http.ErrServerClosed {
		s.Logger.Printf("Server shutted down successfully")
	}

	wg2.Wait()

	return nil
}
