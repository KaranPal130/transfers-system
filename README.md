# Transfers System

This project is an internal transfers system built in Go, designed as part of the assessment. It provides a RESTful API for account management and money transfers, backed by a PostgreSQL database. The system is production-ready, well-documented, and includes interactive Swagger API docs.

## Features
- **Account Creation**: Create new accounts with an initial balance.
- **Balance Query**: Retrieve account balance by account ID.
- **Transaction Submission**: Transfer funds between accounts with validation.
- **Swagger Documentation**: Interactive API documentation at `/swagger/index.html`.
- **Error Handling**: Clear error responses for invalid input, insufficient funds, and more.

## Tech Stack
- **Language**: Go
- **Web Framework**: Gin
- **Database**: PostgreSQL (hosted or local)
- **ORM/SQL**: Standard library + shopspring/decimal
- **API Docs**: swaggo/swag, gin-swagger

## Getting Started

### 1. Prerequisites
- Go 1.18+
- PostgreSQL (local or cloud, connection string in `.env`)

### 2. Setup
1. **Clone the repository**
   ```sh
   git clone <your-repo-url>
   cd transfers-system
   ```
2. **Configure environment variables**
   - Edit `.env` to set your `DATABASE_URL` (see example in the file).

3. **Install dependencies**
   ```sh
   go mod tidy
   ```

4. **Generate Swagger docs**
   ```sh
   swag init -g cmd/server/main.go --output docs
   ```

5. **Run database migrations**
   - Use the provided `scripts/schema.sql` to create tables:
     ```sh
     psql <your-connection-string> -f scripts/schema.sql
     ```

### 3. Running the Server
```sh
go run ./cmd/server
```
The server will start on the port defined in `.env` (default: 8080).

### 4. API Documentation
- Interactive docs: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## API Endpoints

### Account
- `POST /accounts` – Create a new account
- `GET /accounts/{account_id}` – Get account details

### Transactions
- `POST /transactions` – Submit a transfer between accounts

## Example `.env`
```
DATABASE_URL=postgresql://user:password@host:port/dbname
PORT=8080
```

## Project Structure
```
cmd/server/            # Main entry point
internal/api/          # HTTP handlers and server
internal/services/     # Business logic
internal/repositories/ # Database access
internal/models/       # Data models
scripts/schema.sql     # Database schema
.env                   # Environment variables
```

## Notes
- The project is modular, testable, and follows Go best practices.
- API documentation is always up-to-date with code annotations.
- Error messages are clear and user-friendly.
- Code is ready for review and extension.

---

**Good luck reviewing my assessment! If you have any questions, please reach out.**