package main

import (
	"fmt"
	"log"
	"net/http"

	// Импортируем наши пакеты
	"gateway/internal/api"
	"ledger"

	"github.com/go-chi/chi/v5"
)

func main() {
	ledgerService := ledger.NewLedger()

	handlers := api.NewHandlers(ledgerService)

	r := chi.NewRouter()

	r.Use(api.LoggingMiddleware)

	r.Route("/api", func(r chi.Router) {
		// Маршруты для транзакций
		r.Post("/transactions", handlers.CreateTransaction) // POST /api/transactions
		r.Get("/transactions", handlers.ListTransactions)   // GET /api/transactions

		// Маршруты для бюджетов
		r.Post("/budgets", handlers.CreateOrUpdateBudget) // POST /api/budgets
		r.Get("/budgets", handlers.ListBudgets)           // GET /api/budgets
	})

	// 6. Запускаем сервер.
	port := ":8080"
	fmt.Printf("Server is starting on port %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
