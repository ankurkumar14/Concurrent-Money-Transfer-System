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

## Running Tests

To run the tests:

```
go test ./...
``` 