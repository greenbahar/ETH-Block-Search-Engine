package blockprocessor

import (
	"context"
	"ethereum-tracker-app/cmd/config"
	"ethereum-tracker-app/pkg/customerror"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

type Service interface {
	GetBlockNumber(ctx context.Context) (uint64, error)
	GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	GetBlockByHash(ctx context.Context, blockHash common.Hash) (*types.Block, error)
	GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
	GetTransactionLogs(ctx context.Context, txHash common.Hash) ([]*types.Log, error)
	SubscribeNewHeadersViaWss(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	ExtractEvents(ctx context.Context, txs types.Transactions)

	FetchAndStoreRecentBlocks(ctx context.Context, blockTxChan chan types.Transactions) error
	WokerTransactionProcessor(ctx context.Context, blockTxChan chan types.Transactions, wg *sync.WaitGroup)
	SubscribeToNewGeneratedBlocks(ctx context.Context) error
}

type storageService interface {
	SetBlock(ctx context.Context, block *types.Block) error
	SetTransactionHash(ctx context.Context, blockNumber uint64, txHashHex string) error
	SetLogsByTx(ctx context.Context, txHashHex string, logs []*types.Log) error
	SetLogByAddress(ctx context.Context, addressHex string, log *types.Log) error
	GetLogsByAddress(ctx context.Context, addressHex string) ([]types.Log, error)
}

type ethClient struct {
	config     config.Config
	logger     *log.Logger
	httpClient *ethclient.Client
	wsClient   *rpc.Client
	db         storageService
}

func NewEthClient(ctx context.Context, config config.Config, logger *log.Logger, db storageService) (Service, error) {
	client, err := ethclient.DialContext(ctx, config.EthClientConf.EthereumHttpURL)
	if err != nil {
		return nil, customerror.NewConnectionError("", errors.Wrap(err, "cannot connet to http url of ethereum node"))
	}

	rpcClient, err := rpc.DialContext(ctx, config.EthClientConf.EthereumWSSURL)
	if err != nil {
		return nil, customerror.NewConnectionError("", errors.Wrap(err, "cannot connet to wss url of ethereum node"))
	}

	ethClient := &ethClient{
		config:     config,
		logger:     logger,
		httpClient: client,
		wsClient:   rpcClient,
		db:         db,
	}

	return ethClient, nil
}

// GetBlockNumber retrieves the most recent block's number from the Ethereum blockchain
func (ec *ethClient) GetBlockNumber(ctx context.Context) (uint64, error) {
	latestBlock, err := ec.httpClient.BlockNumber(ctx)
	if err != nil {
		return 0, customerror.NewOnChainDataRetrievalError("", errors.Wrap(err, "cannot get the latest block number of the blockchain"))
	}

	return latestBlock, nil
}

// GetBlockByNumber retrieves a block associated with a specific block number
func (ec *ethClient) GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return ec.httpClient.BlockByNumber(ctx, number)
}

// GetTransactionByHash retrieves a transaction by transaction-hash
func (ec *ethClient) GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	return ec.httpClient.TransactionByHash(ctx, hash)
}

// GetLogs retrieves logs (events) with specific filters, in this task logs of an address
func (ec *ethClient) GetLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return ec.httpClient.FilterLogs(ctx, query)
}

// GetTransactionReceipt retrieves the receipt of a transaction, which contains the logs of the transaction as well
func (ec *ethClient) GetTransactionLogs(ctx context.Context, txHash common.Hash) ([]*types.Log, error) {
	receiptOfTx, err := ec.httpClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, customerror.NewLogRetrievalError("", errors.Wrapf(err, "cannot get the logs of the transaction of hash %v", txHash))
	}

	return receiptOfTx.Logs, nil
}

// SubscribeNewBlocks retrieves new header through http url
func (ec *ethClient) SubscribeNewBlocks(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	return ec.httpClient.SubscribeNewHead(ctx, ch)
}

// SubscribeNewHeadersViaWss retrieves new header through wss API, by which the block-number of newly generated blocks can be retrieved
func (ec *ethClient) SubscribeNewHeadersViaWss(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	return ec.wsClient.EthSubscribe(ctx, ch, "newHeads")
}

// GetBlockByHash retrieves a block by block-hash
func (ec *ethClient) GetBlockByHash(ctx context.Context, blockHash common.Hash) (*types.Block, error) {
	return ec.httpClient.BlockByHash(ctx, blockHash)
}
