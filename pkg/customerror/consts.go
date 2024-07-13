package customerror

import "errors"

type ErrorCode int

const (
	// Client error codes (6xxx)
	ErrCodeInvalidInput ErrorCode = iota + 6001 // 6001
	ErrCodeNotFound                             // 6002
	ErrCodeDatabase                             // 6003
	ErrCodeNetwork                              // 6004
	ErrCodeInternal                             // 6005
)

var (
	// Client error messages
	ErrInvalidInput = errors.New("error invalid input data")
	ErrNotFound     = errors.New("error not found")
	ErrDatabase     = errors.New("error database")
	ErrNetwork      = errors.New("error network")
	ErrInternal     = errors.New("error internal")
)

const (
	// Server error codes (7xxx)
	ErrCodeConnection ErrorCode = iota + 7001 // 7001
	ErrCodeRetrieveOnChainData
	ErrCodeRetrieveBlock
	ErrCodeRetrieveLog
	ErrCodeStorage // 7005
)

var (
	// Server error messages
	ErrConnection          = errors.New("cannot connect to host")
	ErrRetrieveOnChainData = errors.New("cannot retrieve some blockchain data from the node")
	ErrRetrieveLog         = errors.New("cannot retrieve log(s)")
	ErrRetrieveBlock       = errors.New("cannot retrieve block(s)")
	ErrStoreData           = errors.New("cannot set data in data storage")
)
