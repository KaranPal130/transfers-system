package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/KaranPal130/transfers-system/docs"
	"github.com/KaranPal130/transfers-system/internal/api"
	repository "github.com/KaranPal130/transfers-system/internal/repositories"
	service "github.com/KaranPal130/transfers-system/internal/services"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		dbConnStr = "postgres://postgres:postgres@localhost:5432/transfers?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(db, accountRepo, transactionRepo)

	handler := api.NewHandler(accountService, transactionService)

	server := api.NewServer(handler)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("Starting server...")
	if err := server.Start(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
