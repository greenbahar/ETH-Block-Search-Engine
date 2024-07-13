/*
Get and store the block and all transaction hashes in the block
Get and store all events related to each transaction in each block

To keep the storage of 50 blocks updated, in case you want strictly keep only 50 blocks and not more, as a new block comes, then the oldest block will be deleted
This can be done easily by the blocknumber itself and no need to store the data in dubly link-list. Here for simplicity, and also flexibility to extend the project
I assumed that new blocks simply just be added (not deleting oldest block).
*/
package inmemorydb

import (
	"context"
	"ethereum-tracker-app/cmd/config"
	"ethereum-tracker-app/pkg/customerror"
	"fmt"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
)

type Service interface {
	SetBlock(ctx context.Context, block *types.Block) error
	SetTransactionHash(ctx context.Context, blockNumber uint64, txHashHex string) error
	SetLogsByTx(ctx context.Context, txHashHex string, logs []*types.Log) error
	SetLogByAddress(ctx context.Context, addressHex string, log *types.Log) error
	GetLogsByAddress(ctx context.Context, addressHex string) ([]types.Log, error)
}

type inmemoryDB struct {
	config config.Config
	logger *log.Logger

	mu          sync.RWMutex
	blocks      map[uint64]*types.Block
	txHashes    map[uint64]string
	txLogs      map[string][]*types.Log
	addressLogs map[string][]*types.Log
}

func NewInmemortDBService(config config.Config, logger *log.Logger) Service {
	return &inmemoryDB{
		config:      config,
		logger:      logger,
		mu:          sync.RWMutex{},
		blocks:      make(map[uint64]*types.Block),
		txHashes:    make(map[uint64]string),
		txLogs:      make(map[string][]*types.Log),
		addressLogs: make(map[string][]*types.Log),
	}
}

// SetBlock gett and stores the block in database
func (db *inmemoryDB) SetBlock(ctx context.Context, block *types.Block) error {
	db.mu.Lock()
	db.blocks[block.NumberU64()] = block
	db.mu.Unlock()

	return nil
}

// SetTransactio sets transaction hashes in the database
func (db *inmemoryDB) SetTransactionHash(ctx context.Context, blockNumber uint64, txHashHex string) error {
	db.mu.Lock()
	db.txHashes[blockNumber] = txHashHex
	db.mu.Unlock()

	return nil
}

// SetLogsByTx stores all events related to each transaction in each block
func (db *inmemoryDB) SetLogsByTx(ctx context.Context, txHashHex string, logs []*types.Log) error {
	db.mu.Lock()
	db.txLogs[txHashHex] = logs
	db.mu.Unlock()

	return nil
}

// SetLogByAddress stores the log related to an address
func (db *inmemoryDB) SetLogByAddress(ctx context.Context, addressHex string, log *types.Log) error {
	db.mu.RLock()
	logs := db.addressLogs[addressHex]
	db.mu.RUnlock()

	db.mu.Lock()
	logs = append(logs, log)
	db.addressLogs[addressHex] = logs
	db.mu.Unlock()

	return nil
}

// GetLogsByAddress gets all the Logs related to an address
func (db *inmemoryDB) GetLogsByAddress(ctx context.Context, addressHex string) ([]types.Log, error) {
	db.mu.RLock()
	logs, ok := db.addressLogs[addressHex]
	db.mu.RUnlock()
	if !ok {
		return nil, customerror.NewStorageError("log does not exist", fmt.Errorf("no logs exists for address %s", addressHex))
	}

	return returnByValue(logs), nil
}

// purpose: safety. blocking the consumer of above functions to unintentionally modify the datastorage, which in this specific case is a map
func returnByValue[k any](input []*k) []k {
	output := make([]k, len(input))
	for i, v := range input {
		if v != nil {
			output[i] = *v
		}
	}

	return output
}
