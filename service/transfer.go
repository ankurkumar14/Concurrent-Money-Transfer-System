package service

import (
	"strings"
)

// TransferResult represents the result of a transfer operation
type TransferResult struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	From    *Account `json:"from,omitempty"`
	To      *Account `json:"to,omitempty"`
}

// TransferRequest represents a request to transfer money between accounts
type TransferRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

// TransferService handles money transfers between accounts
type TransferService struct {
	accountManager AccountManager
}

// NewTransferService creates a new transfer service with the provided account manager
func NewTransferService(accountManager AccountManager) *TransferService {
	return &TransferService{
		accountManager: accountManager,
	}
}

// Transfer performs a money transfer between two accounts
// To prevent deadlocks, locks are acquired in a consistent order (alphabetically by username)
func (ts *TransferService) Transfer(req TransferRequest) (*TransferResult, error) {
	// Validate request
	if req.Amount <= 0 {
		return &TransferResult{Success: false, Message: ErrInvalidAmount.Error()}, ErrInvalidAmount
	}

	// Check if transferring to same account
	if req.From == req.To {
		return &TransferResult{Success: false, Message: ErrSameAccount.Error()}, ErrSameAccount
	}

	// Get accounts
	fromAccount, err := ts.accountManager.GetAccount(req.From)
	if err != nil {
		return &TransferResult{Success: false, Message: "Source account not found"}, err
	}

	toAccount, err := ts.accountManager.GetAccount(req.To)
	if err != nil {
		return &TransferResult{Success: false, Message: "Destination account not found"}, err
	}

	// To prevent deadlocks, always acquire locks in the same order (by username alphabetically)
	first, second := fromAccount, toAccount
	if strings.Compare(fromAccount.Username, toAccount.Username) > 0 {
		first, second = toAccount, fromAccount
	}

	// Acquire locks in order
	first.Lock()
	defer first.Unlock()
	
	second.Lock()
	defer second.Unlock()

	// Check if source has sufficient funds
	if fromAccount.Balance < req.Amount {
		return &TransferResult{
			Success: false,
			Message: ErrInsufficientFunds.Error(),
		}, ErrInsufficientFunds
	}

	// Perform transfer (no need to use Deposit/Withdraw as we already have the locks)
	fromAccount.Balance -= req.Amount
	toAccount.Balance += req.Amount

	// Prepare success result
	result := &TransferResult{
		Success: true,
		Message: "Transfer completed successfully",
		From:    fromAccount,
		To:      toAccount,
	}

	return result, nil
} 