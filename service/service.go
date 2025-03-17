package service

// AccountManager defines the interface for account management operations
type AccountManager interface {
	// GetAccount retrieves an account by username
	GetAccount(username string) (*Account, error)
	
	// ListAccounts returns all accounts
	ListAccounts() []*Account
	
	// CreateAccount creates a new account with the given username and initial balance
	CreateAccount(username string, initialBalance float64) (*Account, error)
} 