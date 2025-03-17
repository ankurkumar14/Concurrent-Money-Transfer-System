package store

import (
	"errors"
	"sync"

	"money-transfer-system/service"
)

// InMemoryStore represents an in-memory implementation of the account store
type InMemoryStore struct {
	accounts map[string]*service.Account
	mutex    sync.RWMutex
}

// NewInMemoryStore creates a new in-memory store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		accounts: make(map[string]*service.Account),
	}
}

// GetAccount retrieves an account by username
func (s *InMemoryStore) GetAccount(username string) (*service.Account, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	account, exists := s.accounts[username]
	if !exists {
		return nil, service.ErrAccountNotFound
	}

	return account, nil
}

// ListAccounts returns all accounts
func (s *InMemoryStore) ListAccounts() []*service.Account {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	accounts := make([]*service.Account, 0, len(s.accounts))
	for _, acc := range s.accounts {
		accounts = append(accounts, acc)
	}

	return accounts
}

// CreateAccount creates a new account with the given username and initial balance
func (s *InMemoryStore) CreateAccount(username string, initialBalance float64) (*service.Account, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if account already exists
	if _, exists := s.accounts[username]; exists {
		return nil, errors.New("account already exists")
	}

	// Create new account
	account := service.NewAccount(username, initialBalance)
	s.accounts[username] = account

	return account, nil
}

// Setup initializes the store with default accounts
func (s *InMemoryStore) Setup() {
	s.CreateAccount("Mark", 100)
	s.CreateAccount("Jane", 50)
	s.CreateAccount("Adam", 0)
} 
