package api

import (
	"github.com/gorilla/mux"
)

// SetupRoutes configures the routes for the API
func (api *API) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Account routes
	r.HandleFunc("/accounts/{username}", api.GetAccountHandler).Methods("GET")
	r.HandleFunc("/accounts", api.ListAccountsHandler).Methods("GET")

	// Transfer route
	r.HandleFunc("/transfer", api.TransferHandler).Methods("POST")

	return r
} 