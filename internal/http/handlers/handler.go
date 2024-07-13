package handlers

import (
	"encoding/json"
	"errors"
	"ethereum-tracker-app/internal/services/blocksearch"
	"ethereum-tracker-app/models"
	"ethereum-tracker-app/pkg/customerror"
	"net/http"
)

type Handler interface {
	GetEventsByAddress(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	blockProcessService blocksearch.Service
}

func NewHandler(blockProcessorSrv blocksearch.Service) Handler {
	return &handler{
		blockProcessService: blockProcessorSrv,
	}
}

func (h *handler) respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(models.ErrorResponse{Code: code, Message: message})
}

func (h *handler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func (h *handler) handleError(w http.ResponseWriter, err error) {
	var customErr *customerror.Error
	if errors.As(err, &customErr) {
		switch customErr.Code {
		case customerror.ErrCodeNotFound:
			h.respondWithError(w, http.StatusNotFound, customErr.Message)
		case customerror.ErrCodeInvalidInput:
			h.respondWithError(w, http.StatusBadRequest, customErr.Message)
		default:
			h.respondWithError(w, http.StatusInternalServerError, customErr.Message)
		}

		return
	}

	h.respondWithError(w, http.StatusInternalServerError, err.Error())
}
