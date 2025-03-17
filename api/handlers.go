package api

import (
	"encoding/json"
	"net/http"

	"money-transfer-system/service"

	"github.com/gorilla/mux"
)

// API holds the handlers and services for the HTTP API
type API struct {
	transferService *service.TransferService
	accountManager  service.AccountManager
}

// NewAPI creates a new API instance
func NewAPI(transferService *service.TransferService, accountManager service.AccountManager) *API {
	return &API{
		transferService: transferService,
		accountManager:  accountManager,
	}
}

// GetAccountHandler returns the account information for the specified username
func (api *API) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	account, err := api.accountManager.GetAccount(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// ListAccountsHandler returns a list of all accounts
func (api *API) ListAccountsHandler(w http.ResponseWriter, r *http.Request) {
	accounts := api.accountManager.ListAccounts()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

// TransferHandler handles money transfer requests
func (api *API) TransferHandler(w http.ResponseWriter, r *http.Request) {
	var req service.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	result, err := api.transferService.Transfer(req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
} 
