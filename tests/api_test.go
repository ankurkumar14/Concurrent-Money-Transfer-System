package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"money-transfer-system/api"
	"money-transfer-system/service"
	"money-transfer-system/store"
)

func setupTestAPI() *api.API {
	accountStore := store.NewInMemoryStore()
	accountStore.CreateAccount("Mark", 100)
	accountStore.CreateAccount("Jane", 50)
	accountStore.CreateAccount("Adam", 0)

	transferService := service.NewTransferService(accountStore)
	apiHandler := api.NewAPI(transferService, accountStore)

	return apiHandler
}

func TestGetAccountHandler(t *testing.T) {
	// Setup
	apiHandler := setupTestAPI()
	router := apiHandler.SetupRoutes()

	// Create a request to get Mark's account
	req, _ := http.NewRequest("GET", "/accounts/Mark", nil)
	rr := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", rr.Code)
	}

	// Parse response
	var account service.Account
	if err := json.Unmarshal(rr.Body.Bytes(), &account); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	// Verify account data
	if account.Username != "Mark" {
		t.Errorf("Expected username Mark, got %v", account.Username)
	}

	if account.Balance != 100 {
		t.Errorf("Expected balance 100, got %v", account.Balance)
	}
}

func TestListAccountsHandler(t *testing.T) {
	// Setup
	apiHandler := setupTestAPI()
	router := apiHandler.SetupRoutes()

	// Create a request to list all accounts
	req, _ := http.NewRequest("GET", "/accounts", nil)
	rr := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", rr.Code)
	}

	// Parse response
	var accounts []*service.Account
	if err := json.Unmarshal(rr.Body.Bytes(), &accounts); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	// Verify accounts count
	if len(accounts) != 3 {
		t.Errorf("Expected 3 accounts, got %v", len(accounts))
	}
}

func TestTransferHandler(t *testing.T) {
	// Setup
	apiHandler := setupTestAPI()
	router := apiHandler.SetupRoutes()

	// Create a transfer request (Mark sends $25 to Jane)
	transferReq := service.TransferRequest{
		From:   "Mark",
		To:     "Jane",
		Amount: 25,
	}
	reqBody, _ := json.Marshal(transferReq)

	// Create HTTP request
	req, _ := http.NewRequest("POST", "/transfer", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", rr.Code)
	}

	// Parse response
	var result service.TransferResult
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	// Verify transfer was successful
	if !result.Success {
		t.Errorf("Expected transfer success, got failure: %v", result.Message)
	}

	// Verify balances
	if result.From.Balance != 75 {
		t.Errorf("Expected sender balance 75, got %v", result.From.Balance)
	}

	if result.To.Balance != 75 {
		t.Errorf("Expected recipient balance 75, got %v", result.To.Balance)
	}
}

func TestTransferInsufficientFunds(t *testing.T) {
	// Setup
	apiHandler := setupTestAPI()
	router := apiHandler.SetupRoutes()

	// Create a transfer request (Adam tries to send $50 to Jane, but has $0)
	transferReq := service.TransferRequest{
		From:   "Adam",
		To:     "Jane",
		Amount: 50,
	}
	reqBody, _ := json.Marshal(transferReq)

	// Create HTTP request
	req, _ := http.NewRequest("POST", "/transfer", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rr.Code)
	}

	// Parse response
	var result service.TransferResult
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	// Verify transfer failed due to insufficient funds
	if result.Success {
		t.Errorf("Expected transfer failure, got success")
	}

	if result.Message != service.ErrInsufficientFunds.Error() {
		t.Errorf("Expected error message '%v', got '%v'", service.ErrInsufficientFunds.Error(), result.Message)
	}
} 
