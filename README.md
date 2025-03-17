# Concurrent Money Transfer System

A simple concurrent money transfer system built in Go that allows users to transfer money between accounts while maintaining thread safety and preventing overdrafts.

## Features

- Concurrent money transfers between users
- Thread-safe account operations using mutex locks
- Overdraft prevention
- HTTP API for initiating transfers
- Initial balances: Mark ($100), Jane ($50), Adam ($0)

## Setup and Installation

### Prerequisites

- Go 1.16 or later

### Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd money-transfer-system
   ```

2. Build the application:
   ```
   go build -o transfer-app
   ```

3. Run the application:
   ```
   ./transfer-app
   ```

The server will start on port 8080 by default.

## API Documentation

### Get Account Balance

```
GET /accounts/{username}
```

Returns the current balance of the specified user.

**Response:**
```json
{
  "username": "Mark",
  "balance": 100
}
```

### List All Accounts

```
GET /accounts
```

Returns a list of all accounts and their balances.

**Response:**
```json
[
  {
    "username": "Mark",
    "balance": 100
  },
  {
    "username": "Jane",
    "balance": 50
  },
  {
    "username": "Adam",
    "balance": 0
  }
]
```

### Transfer Money

```
POST /transfer
```

Initiates a money transfer between two users.

**Request Body:**
```json
{
  "from": "Mark",
  "to": "Jane",
  "amount": 25
}
```

**Success Response:**
```json
{
  "success": true,
  "message": "Transfer completed successfully",
  "from": {
    "username": "Mark",
    "balance": 75
  },
  "to": {
    "username": "Jane",
    "balance": 75
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Insufficient funds"
}
```

## Concurrency Strategy

The system uses a combination of account-level mutex locks to ensure atomic operations when updating account balances. For transfers, we acquire locks on both source and destination accounts to prevent race conditions.

To avoid deadlocks, the system always acquires locks in the same order (based on account username) regardless of the transfer direction.
The system is designed to handle multiple concurrent transfers safely.
1. Account-Level Locking: Each account has its own mutex lock (sync.Mutex) that prevents concurrent access to the account's balance. This ensures that operations on a single account (like deposits and withdrawals) are atomic.
2. Deadlock Prevention: The transfer service implements a critical deadlock prevention strategy by always acquiring locks in a consistent order (alphabetically by username). This is a standard solution to the dining philosophers problem and prevents circular wait conditions.
3. Atomic Transfers: The entire transfer operation (checking balances, deducting from source, adding to destination) is performed as an atomic operation while holding locks on both accounts. This ensures consistency and prevents race conditions.
4. Test Results: Our tests confirm that the system handles concurrency correctly:
The original TestConcurrentTransfers test passed, showing that 10 concurrent transfers complete correctly.
Our custom TestHighConcurrencyTransfers test with 500 concurrent transfers (100 iterations × 5 transfers) passed, maintaining the correct total balance.
Our TestRandomConcurrentTransfers test with 500 random transfers between 10 accounts also passed, showing that the system can handle a more realistic workload with varying transfer amounts and directions.

##Concurrency Mechanisms Used
1. Mutex Locks: Each account has a mutex that prevents concurrent access to its balance.
2. Consistent Lock Ordering: Locks are always acquired in the same order to prevent deadlocks.
3. Deferred Unlocking: The defer statement ensures locks are released even if errors occur.
4. Atomic Operations: Balance checks and updates are performed atomically while holding locks.

##Potential Concurrency Issues Addressed
1. Race Conditions: Prevented by using mutex locks to ensure exclusive access to account balances.
2. Deadlocks: Prevented by acquiring locks in a consistent order.
3. Lost Updates: Prevented by making transfers atomic operations.
4. Inconsistent State: Prevented by checking balances while holding locks.

## Running Tests

To run the tests:

```
go test ./...
```

## Conclusion
The money transfer system is well-designed to handle concurrent transfers from multiple users simultaneously. The implementation uses proper synchronization techniques to ensure thread safety, prevent race conditions, and maintain data consistency.
Our tests demonstrate that the system can handle hundreds of concurrent transfers without issues, maintaining the correct total balance across all accounts. This indicates that the concurrency mechanisms are working as intended.
In a real-world scenario, this system would be able to handle many users making transfers at the same time, which is essential for a financial application where multiple transactions need to be processed concurrently.
