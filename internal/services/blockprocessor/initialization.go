package blockprocessor

import (
	"context"
	"ethereum-tracker-app/pkg/customerror"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

// FetchAndStoreRecentBlocks retrieves blocks from node and store the data in storage
func (ec *ethClient) FetchAndStoreRecentBlocks(ctx context.Context, blockTxChan chan types.Transactions) error {
	latestBlock, err := ec.GetBlockNumber(ctx)
	if err != nil {
		return customerror.NewBlockRetrievalError("", errors.Wrap(err, "cannot fetch the most latest block number"))
	}

	for i := uint64(0); i < uint64(ec.config.EthClientConf.NumberOfRecentBlocks); i++ {
		select {
		case <-ctx.Done():
			close(blockTxChan)
			ec.logger.Println("Context cancelled, stopping FetchAndStoreRecentBlocks processor")
			return nil
		default:
			blockNumber := latestBlock - i
			if blockNumber == 0 {
				break
			}
			ec.logger.Printf("block %d recieved", blockNumber)

			block, err := ec.GetBlockByNumber(ctx, big.NewInt(int64(blockNumber)))
			if err != nil {
				ec.logger.Printf("cannot retreive block %d", blockNumber)
				continue
			}

			if setErr := ec.db.SetBlock(ctx, block); setErr != nil {
				ec.logger.Printf("cannot set the block %d", blockNumber)
			}

			blockTxChan <- block.Transactions()
		}
	}

	// not closing the chanels is a common cause of the goroutine leak as they never stop
	close(blockTxChan)

	return nil
}

// WokerTransactionProcessor is a worker to process the tranactions of a block
func (ec *ethClient) WokerTransactionProcessor(ctx context.Context, blockTxChan chan types.Transactions, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			ec.logger.Println("Context cancelled, stopping transaction processor")
			return
		case txs, ok := <-blockTxChan:
			if !ok {
				// Channel closed, exit the loop
				ec.logger.Println("blockTxChan closed, stopping transaction processor")
				return
			}

			ec.ExtractEvents(ctx, txs)
		}
	}
}

// ExtractEvents gets the events (logs) of transactions in a block
func (ec *ethClient) ExtractEvents(ctx context.Context, txs types.Transactions) {
	for _, tx := range txs {
		logs, err := ec.GetTransactionLogs(ctx, tx.Hash())
		if err != nil {
			ec.logger.Printf("failed to get logs of the transaction of hash %v \n", tx.Hash())
		}

		if len(logs) != 0 {
			if setLogErr := ec.db.SetLogsByTx(ctx, tx.Hash().Hex(), logs); setLogErr != nil {
				ec.logger.Printf("failed to store logs of transaction %v \n", tx.Hash())
			}

			for _, txLog := range logs {
				if logAdrErr := ec.db.SetLogByAddress(ctx, txLog.Address.Hex(), txLog); logAdrErr != nil {
					ec.logger.Printf("faled to store log %v in transaction %v in block %d related to address %v  \n", txLog, tx.Hash(), txLog.BlockNumber, txLog.Address)
				}
			}
		}
	}
}
