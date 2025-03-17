package tests

import (
	"sync"
	"testing"

	"money-transfer-system/service"
	"money-transfer-system/store"
)

func TestHighConcurrencyTransfers(t *testing.T) {
	// Setup
	accountStore := store.NewInMemoryStore()
	
	// Create 5 accounts with initial balance of 1000 each
	accountStore.CreateAccount("User1", 1000)
	accountStore.CreateAccount("User2", 1000)
	accountStore.CreateAccount("User3", 1000)
	accountStore.CreateAccount("User4", 1000)
	accountStore.CreateAccount("User5", 1000)
	
	transferService := service.NewTransferService(accountStore)
	
	// Use WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	
	// Number of concurrent transfers
	numTransfers := 100
	
	// Execute many concurrent transfers between all accounts
	for i := 0; i < numTransfers; i++ {
		wg.Add(5) // 5 transfers per iteration
		
		// User1 -> User2
		go func() {
			defer wg.Done()
			req := service.TransferRequest{From: "User1", To: "User2", Amount: 1}
			transferService.Transfer(req)
		}()
		
		// User2 -> User3
		go func() {
			defer wg.Done()
			req := service.TransferRequest{From: "User2", To: "User3", Amount: 1}
			transferService.Transfer(req)
		}()
		
		// User3 -> User4
		go func() {
			defer wg.Done()
			req := service.TransferRequest{From: "User3", To: "User4", Amount: 1}
			transferService.Transfer(req)
		}()
		
		// User4 -> User5
		go func() {
			defer wg.Done()
			req := service.TransferRequest{From: "User4", To: "User5", Amount: 1}
			transferService.Transfer(req)
		}()
		
		// User5 -> User1
		go func() {
			defer wg.Done()
			req := service.TransferRequest{From: "User5", To: "User1", Amount: 1}
			transferService.Transfer(req)
		}()
	}
	
	// Wait for all transfers to complete
	wg.Wait()
	
	// Get final balances
	user1, _ := accountStore.GetAccount("User1")
	user2, _ := accountStore.GetAccount("User2")
	user3, _ := accountStore.GetAccount("User3")
	user4, _ := accountStore.GetAccount("User4")
	user5, _ := accountStore.GetAccount("User5")
	
	// Calculate total balance (should remain constant)
	totalBalance := user1.GetBalance() + user2.GetBalance() + user3.GetBalance() + 
		user4.GetBalance() + user5.GetBalance()
	
	// Verify total balance is still 5000 (1000 * 5)
	if totalBalance != 5000 {
		t.Errorf("Expected total balance=5000, got %v", totalBalance)
	}
	
	// Log final balances
	t.Logf("Final balances after %d concurrent transfers per account:", numTransfers)
	t.Logf("User1: %.2f", user1.GetBalance())
	t.Logf("User2: %.2f", user2.GetBalance())
	t.Logf("User3: %.2f", user3.GetBalance())
	t.Logf("User4: %.2f", user4.GetBalance())
	t.Logf("User5: %.2f", user5.GetBalance())
} 