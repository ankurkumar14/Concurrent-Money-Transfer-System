package tests

import (
	"testing"

	"money-transfer-system/service"
	"money-transfer-system/store"
)

func TestSuccessfulTransfer(t *testing.T) {
	// Setup
	accountStore := store.NewInMemoryStore()
	accountStore.CreateAccount("User1", 100)
	accountStore.CreateAccount("User2", 50)
	
	transferService := service.NewTransferService(accountStore)
	
	// Execute transfer
	req := service.TransferRequest{
		From:   "User1",
		To:     "User2",
		Amount: 25,
	}
	
	result, err := transferService.Transfer(req)
	
	// Verify
	if err != nil {
		t.Errorf("Expected successful transfer, got error: %v", err)
	}
	
	if !result.Success {
		t.Errorf("Expected success=true, got success=%v", result.Success)
	}
	
	// Check balances
	user1, _ := accountStore.GetAccount("User1")
	user2, _ := accountStore.GetAccount("User2")
	
	if user1.GetBalance() != 75 {
		t.Errorf("Expected User1 balance=75, got %v", user1.GetBalance())
	}
	
	if user2.GetBalance() != 75 {
		t.Errorf("Expected User2 balance=75, got %v", user2.GetBalance())
	}
}

func TestInsufficientFunds(t *testing.T) {
	// Setup
	accountStore := store.NewInMemoryStore()
	accountStore.CreateAccount("User1", 100)
	accountStore.CreateAccount("User2", 50)
	
	transferService := service.NewTransferService(accountStore)
	
	// Execute transfer with insufficient funds
	req := service.TransferRequest{
		From:   "User1",
		To:     "User2",
		Amount: 150, // More than available
	}
	
	result, err := transferService.Transfer(req)
	
	// Verify
	if err != service.ErrInsufficientFunds {
		t.Errorf("Expected ErrInsufficientFunds, got: %v", err)
	}
	
	if result.Success {
		t.Errorf("Expected success=false, got success=%v", result.Success)
	}
	
	// Check balances are unchanged
	user1, _ := accountStore.GetAccount("User1")
	user2, _ := accountStore.GetAccount("User2")
	
	if user1.GetBalance() != 100 {
		t.Errorf("Expected User1 balance=100, got %v", user1.GetBalance())
	}
	
	if user2.GetBalance() != 50 {
		t.Errorf("Expected User2 balance=50, got %v", user2.GetBalance())
	}
}

func TestTransferToSelf(t *testing.T) {
	// Setup
	accountStore := store.NewInMemoryStore()
	accountStore.CreateAccount("User1", 100)
	
	transferService := service.NewTransferService(accountStore)
	
	// Execute transfer to self
	req := service.TransferRequest{
		From:   "User1",
		To:     "User1",
		Amount: 25,
	}
	
	result, err := transferService.Transfer(req)
	
	// Verify
	if err != service.ErrSameAccount {
		t.Errorf("Expected ErrSameAccount, got: %v", err)
	}
	
	if result.Success {
		t.Errorf("Expected success=false, got success=%v", result.Success)
	}
	
	// Check balance is unchanged
	user1, _ := accountStore.GetAccount("User1")
	if user1.GetBalance() != 100 {
		t.Errorf("Expected User1 balance=100, got %v", user1.GetBalance())
	}
}

func TestInvalidAmount(t *testing.T) {
	// Setup
	accountStore := store.NewInMemoryStore()
	accountStore.CreateAccount("User1", 100)
	accountStore.CreateAccount("User2", 50)
	
	transferService := service.NewTransferService(accountStore)
	
	// Execute transfer with negative amount
	req := service.TransferRequest{
		From:   "User1",
		To:     "User2",
		Amount: -25,
	}
	
	result, err := transferService.Transfer(req)
	
	// Verify
	if err != service.ErrInvalidAmount {
		t.Errorf("Expected ErrInvalidAmount, got: %v", err)
	}
	
	if result.Success {
		t.Errorf("Expected success=false, got success=%v", result.Success)
	}
	
	// Check balances are unchanged
	user1, _ := accountStore.GetAccount("User1")
	user2, _ := accountStore.GetAccount("User2")
	
	if user1.GetBalance() != 100 {
		t.Errorf("Expected User1 balance=100, got %v", user1.GetBalance())
	}
	
	if user2.GetBalance() != 50 {
		t.Errorf("Expected User2 balance=50, got %v", user2.GetBalance())
	}
}

func TestConcurrentTransfers(t *testing.T) {
	// Setup
	accountStore := store.NewInMemoryStore()
	accountStore.CreateAccount("User1", 100)
	accountStore.CreateAccount("User2", 100)
	
	transferService := service.NewTransferService(accountStore)
	
	// Create a channel to synchronize goroutines
	done := make(chan bool)
	
	// Execute 5 concurrent transfers back and forth
	for i := 0; i < 5; i++ {
		go func() {
			req1 := service.TransferRequest{From: "User1", To: "User2", Amount: 10}
			transferService.Transfer(req1)
			done <- true
		}()
		
		go func() {
			req2 := service.TransferRequest{From: "User2", To: "User1", Amount: 5}
			transferService.Transfer(req2)
			done <- true
		}()
	}
	
	// Wait for all transfers to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Verify final balances
	// Should be: User1 = 100 - (5*10) + (5*5) = 75
	//            User2 = 100 + (5*10) - (5*5) = 125
	user1, _ := accountStore.GetAccount("User1")
	user2, _ := accountStore.GetAccount("User2")
	
	if user1.GetBalance() != 75 {
		t.Errorf("Expected User1 balance=75, got %v", user1.GetBalance())
	}
	
	if user2.GetBalance() != 125 {
		t.Errorf("Expected User2 balance=125, got %v", user2.GetBalance())
	}
} 