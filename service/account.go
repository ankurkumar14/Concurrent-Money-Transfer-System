package service

import (
	"errors"
	"sync"
)

// Common errors
var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrAccountNotFound   = errors.New("account not found")
	ErrInvalidAmount     = errors.New("invalid amount, must be positive")
	ErrSameAccount       = errors.New("cannot transfer to the same account")
)

// Account represents a user account with balance
type Account struct {
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
	mutex    sync.Mutex
}

// NewAccount creates a new account with the given username and initial balance
func NewAccount(username string, initialBalance float64) *Account {
	return &Account{
		Username: username,
		Balance:  initialBalance,
	}
}

// Deposit adds the specified amount to the account balance
// Returns an error if the amount is negative
func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.Balance += amount
	return nil
}

// Withdraw subtracts the specified amount from the account balance
// Returns an error if there are insufficient funds or if the amount is negative
func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.Balance < amount {
		return ErrInsufficientFunds
	}

	a.Balance -= amount
	return nil
}

// GetBalance returns the current balance of the account
func (a *Account) GetBalance() float64 {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	return a.Balance
}

// Lock locks the account for concurrent access
func (a *Account) Lock() {
	a.mutex.Lock()
}

// Unlock unlocks the account
func (a *Account) Unlock() {
	a.mutex.Unlock()
} 
