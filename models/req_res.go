package models

import (
	"github.com/ethereum/go-ethereum/core/types"
)

// ErrorResponse represents a standard error response. The code in the response is an http status code and not the internal service error codes
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// EventResponse represents the successful response containing events
type EventResponse struct {
	Status  Status      `json:"status"`
	Address string      `json:"address"`
	Events  []types.Log `json:"events"`
}

type Status string

const (
	StatusSuccess Status = "success"
	StatusCreated Status = "created"
)
