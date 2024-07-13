package handlers

import (
	"ethereum-tracker-app/models"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
)

// Events API endpoint
// @Summary Get events by address
// @Description Retrieve events (very lgs) of a specific address
// @Tags Logs
// @Accept json
// @Produce json
// @Param address path string true "an address in the blockchain"
// @Success 200 {array} types.Log
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /events/{address} [get]

// GetEventsByAddress Gets the events related to a specific address
func (h *handler) GetEventsByAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	if !common.IsHexAddress(address) {
		h.respondWithError(w, http.StatusBadRequest, "Invalid input: address is not a valid hex address")
		return
	}

	events, err := h.blockProcessService.GetEventsByAddress(r.Context(), address)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, models.EventResponse{
		Status:  models.StatusSuccess,
		Address: address,
		Events:  events,
	})
}
