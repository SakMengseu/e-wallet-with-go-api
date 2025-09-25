# e-Wallet API

A simple e-Wallet REST API built with **Go (Gin)** and **MongoDB**, allowing users to register, login, manage wallets, and perform transactions.

---

## Table of Contents

* [Features](#features)
* [Tech Stack](#tech-stack)
* [API Routes](#api-routes)

  * [Authentication](#authentication)
  * [Wallet](#wallet)
  * [Transactions](#transactions)
* [Setup & Installation](#setup--installation)
* [Usage](#usage)
* [License](#license)

---

## Features

* User registration and login with JWT authentication
* Wallet management: get balance, deposit, withdraw
* Money transfers between users
* Transaction history and details
* Secure endpoints using middleware authentication

---

## Tech Stack

* **Go** with **Gin framework**
* **MongoDB**
* JWT for authentication
* Go modules for dependency management

---

## API Routes

### Authentication

| Method | Endpoint           | Description         | Body                                                                      |
| ------ | ------------------ | ------------------- | ------------------------------------------------------------------------- |
| POST   | `/api/v1/register` | Register a new user | `{ "name": "Alice", "email": "alice@example.com", "password": "123456" }` |
| POST   | `/api/v1/login`    | Login and get JWT   | `{ "email": "alice@example.com", "password": "123456" }`                  |

---

### Wallet

All wallet routes require **Authorization header**:

```
Authorization: Bearer <JWT_TOKEN>
```

| Method | Endpoint                  | Description                | Body                |
| ------ | ------------------------- | -------------------------- | ------------------- |
| GET    | `/api/v1/wallet`          | Get current wallet balance | -                   |
| POST   | `/api/v1/wallet/deposit`  | Deposit money into wallet  | `{ "amount": 100 }` |
| POST   | `/api/v1/wallet/withdraw` | Withdraw money from wallet | `{ "amount": 50 }`  |

---

### Transactions

| Method | Endpoint                                       | Description                                 | Body                                         |
| ------ | ---------------------------------------------- | ------------------------------------------- | -------------------------------------------- |
| POST   | `/api/v1/transactions/send`                    | Send money to another user                  | `{ "receiver_id": "USER_ID", "amount": 50 }` |
| GET    | `/api/v1/transactions/histories`               | Get all transactions for the logged-in user | -                                            |
| GET    | `/api/v1/transactions/history/:transaction_id` | Get details of a specific transaction       | -                                            |

---

## Setup & Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/e-wallet.git
cd e-wallet
```

2. Install dependencies:

```bash
go mod tidy
```

3. Set up **MongoDB** and update the connection string in your config.

4. Run the application:

```bash
go run main.go
```

5. The API will be available at `http://localhost:8080`.

---

## Usage

1. **Register a user** via `/api/v1/register`
2. **Login** to get JWT token via `/api/v1/login`
3. Include the token in the header for all protected routes:

```
Authorization: Bearer <JWT_TOKEN>
```

4. Perform **wallet operations** or **transactions** as needed.

---

## License

SAK MENGSEU License © 2025
