/*
	Expose an endpoint that allows a user to query all "events" related to a particular "address"
*/

package routers

import (
	_ "ethereum-tracker-app/docs"
	"ethereum-tracker-app/internal/http/handlers"
	"net/http"

	_ "github.com/ethereum/go-ethereum/core/types"

	"github.com/gorilla/mux"
)

// @title Ethereum Tracker API
// @version 1.0
// @description API endpoints for Ethereum blockchain tracking
// @basePath /v1
func SetupRouters(handler handlers.Handler) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/v1/events/{address}", handler.GetEventsByAddress).Methods("GET")

	// Serve the Swagger documentation JSON
	router.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../docs/swagger.json")
	})

	return router
}
