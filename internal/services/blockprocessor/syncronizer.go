package blockprocessor

import (
	"context"
	"ethereum-tracker-app/pkg/customerror"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

// SubscribeToNewGeneratedBlocks retrieves and stores new generated blocks
func (ec *ethClient) SubscribeToNewGeneratedBlocks(ctx context.Context) error {
	headers := make(chan *types.Header)
	sub, err := ec.SubscribeNewHeadersViaWss(ctx, headers)
	if err != nil {
		return customerror.NewOnChainDataRetrievalError("", errors.Wrapf(err, "failed to subscribe to new headers of the blockchain"))
	}

	for {
		select {
		case <-ctx.Done():
			ec.logger.Println("Context cancelled, stopping block subscription")
			return nil
		case err := <-sub.Err():
			ec.logger.Printf(customerror.NewOnChainDataRetrievalError("error in header subscription", err).Error())
		case header := <-headers:
			ec.logger.Printf("New block received: %v \n", header.Number.String())

			block, err := ec.GetBlockByNumber(ctx, header.Number)
			if err != nil {
				ec.logger.Printf("failed to fetch block details of block %d: %v", header.Number, err)
				continue
			}

			if setErr := ec.db.SetBlock(context.Background(), block); setErr != nil {
				ec.logger.Printf("block %d has not been stored in the datastore", block.NumberU64())
				// todo having exra mechanism to handle this occasion to store blocks in case of error
			}

			ec.ExtractEvents(ctx, block.Transactions())
		}
	}
}
