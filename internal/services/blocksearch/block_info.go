package blocksearch

import (
	"context"
	"ethereum-tracker-app/cmd/config"
	"ethereum-tracker-app/pkg/customerror"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type Service interface {
	GetEventsByAddress(ctx context.Context, address string) ([]types.Log, error)
}

type storageService interface {
	SetBlock(ctx context.Context, block *types.Block) error
	SetTransactionHash(ctx context.Context, blockNumber uint64, txHashHex string) error
	SetLogsByTx(ctx context.Context, txHashHex string, logs []*types.Log) error
	SetLogByAddress(ctx context.Context, addressHex string, log *types.Log) error
	GetLogsByAddress(ctx context.Context, addressHex string) ([]types.Log, error)
}

type blockprocess struct {
	config config.Config
	logger *log.Logger
	db     storageService
}

func NewServie(config config.Config, logger *log.Logger, db storageService) Service {
	return &blockprocess{
		config: config,
		logger: logger,
		db:     db,
	}
}

// GetEventsByAddress gets events of a specific address
func (b *blockprocess) GetEventsByAddress(ctx context.Context, address string) ([]types.Log, error) {
	logs, err := b.db.GetLogsByAddress(ctx, address)
	if err != nil {
		return nil, customerror.NewLogRetrievalError("", errors.Wrapf(err, "failed to get logs of address %s", address))
	}

	return logs, nil
}
