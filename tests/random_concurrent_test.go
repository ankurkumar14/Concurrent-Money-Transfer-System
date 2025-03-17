package tests

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"money-transfer-system/service"
	"money-transfer-system/store"
)

func TestRandomConcurrentTransfers(t *testing.T) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Setup
	accountStore := store.NewInMemoryStore()
	
	// Create 10 accounts with initial balance of 1000 each
	usernames := []string{
		"Alice", "Bob", "Charlie", "Dave", "Eve",
		"Frank", "Grace", "Heidi", "Ivan", "Judy",
	}
	
	initialBalance := 1000.0
	initialTotalBalance := initialBalance * float64(len(usernames))
	
	for _, name := range usernames {
		accountStore.CreateAccount(name, initialBalance)
	}
	
	transferService := service.NewTransferService(accountStore)
	
	// Use WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	
	// Number of concurrent transfers
	numTransfers := 500
	wg.Add(numTransfers)
	
	// Track successful and failed transfers
	var successCount, failCount int
	var countMutex sync.Mutex
	
	// Execute random concurrent transfers
	for i := 0; i < numTransfers; i++ {
		go func() {
			defer wg.Done()
			
			// Select random source and destination accounts
			fromIndex := rand.Intn(len(usernames))
			toIndex := rand.Intn(len(usernames))
			
			// Skip if same account
			if fromIndex == toIndex {
				countMutex.Lock()
				failCount++
				countMutex.Unlock()
				return
			}
			
			// Random amount between 1 and 100
			amount := float64(rand.Intn(100) + 1)
			
			// Perform transfer
			req := service.TransferRequest{
				From:   usernames[fromIndex],
				To:     usernames[toIndex],
				Amount: amount,
			}
			
			result, _ := transferService.Transfer(req)
			
			// Count successful and failed transfers
			countMutex.Lock()
			if result.Success {
				successCount++
			} else {
				failCount++
			}
			countMutex.Unlock()
		}()
	}
	
	// Wait for all transfers to complete
	wg.Wait()
	
	// Calculate final total balance
	var finalTotalBalance float64
	for _, name := range usernames {
		account, _ := accountStore.GetAccount(name)
		finalTotalBalance += account.GetBalance()
		
		// Log individual balances
		t.Logf("%s final balance: %.2f", name, account.GetBalance())
	}
	
	// Verify total balance is preserved
	if finalTotalBalance != initialTotalBalance {
		t.Errorf("Expected total balance=%.2f, got %.2f", initialTotalBalance, finalTotalBalance)
	}
	
	// Log transfer statistics
	t.Logf("Completed %d transfers: %d successful, %d failed", numTransfers, successCount, failCount)
} 