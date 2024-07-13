package customerror

import (
	"fmt"
)

type Error struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %d, message: %s, error: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func NewConnectionError(message string, err error) *Error {
	if message == "" && err != nil {
		return New(ErrCodeConnection, ErrConnection.Error(), err)
	}
	if err == nil && message != "" {
		return New(ErrCodeConnection, message, ErrConnection)
	}

	return New(ErrCodeConnection, message, err)
}

func NewOnChainDataRetrievalError(message string, err error) *Error {
	if message == "" && err != nil {
		return New(ErrCodeRetrieveOnChainData, ErrRetrieveOnChainData.Error(), err)
	}
	if err == nil && message != "" {
		return New(ErrCodeRetrieveOnChainData, message, ErrRetrieveOnChainData)
	}

	return New(ErrCodeRetrieveOnChainData, message, err)
}

func NewBlockRetrievalError(message string, err error) *Error {
	if message == "" && err != nil {
		return New(ErrCodeRetrieveBlock, ErrRetrieveBlock.Error(), err)
	}
	if err == nil && message != "" {
		return New(ErrCodeRetrieveBlock, message, ErrRetrieveBlock)
	}

	return New(ErrCodeRetrieveBlock, message, err)
}

func NewLogRetrievalError(message string, err error) *Error {
	if message == "" && err != nil {
		return New(ErrCodeRetrieveLog, ErrRetrieveLog.Error(), err)
	}
	if err == nil && message != "" {
		return New(ErrCodeRetrieveLog, message, ErrRetrieveLog)
	}

	return New(ErrCodeRetrieveLog, message, err)
}

func NewStorageError(message string, err error) *Error {
	if message == "" && err != nil {
		return New(ErrCodeStorage, ErrDatabase.Error(), err)
	}
	if err == nil && message != "" {
		return New(ErrCodeStorage, message, ErrDatabase)
	}

	return New(ErrCodeStorage, message, err)
}
