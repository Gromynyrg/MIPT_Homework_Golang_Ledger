package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ledger"

	"github.com/go-chi/chi/v5"
)

// setupTestServer - вспомогательная функция для создания тестового сервера.
func setupTestServer() (*httptest.Server, *ledger.Ledger) {
	ledgerService := ledger.NewLedger()
	handlers := NewHandlers(ledgerService)
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Post("/transactions", handlers.CreateTransaction)
		r.Get("/transactions", handlers.ListTransactions)
		r.Post("/budgets", handlers.CreateOrUpdateBudget)
		r.Get("/budgets", handlers.ListBudgets)
	})

	server := httptest.NewServer(r)
	return server, ledgerService
}

func TestBudgetsAPI(t *testing.T) {
	server, ledgerService := setupTestServer()
	defer server.Close()

	t.Cleanup(func() {
		ledgerService.Reset()
	})

	t.Run("Create and Get Budget", func(t *testing.T) {
		budgetPayload := `{"category":"test","limit":100}`
		resp, err := http.Post(server.URL+"/api/budgets", "application/json", bytes.NewBufferString(budgetPayload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201 Created, got %d", resp.StatusCode)
		}

		resp, err = http.Get(server.URL + "/api/budgets")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
		}

		var budgets []BudgetResponse
		if err := json.NewDecoder(resp.Body).Decode(&budgets); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(budgets) != 1 || budgets[0].Category != "test" || budgets[0].Limit != 100 {
			t.Errorf("Unexpected budgets list: %+v", budgets)
		}
	})

	t.Run("Invalid budget payload", func(t *testing.T) {
		payload := `{"category":"test","limit":-100}`
		resp, err := http.Post(server.URL+"/api/budgets", "application/json", bytes.NewBufferString(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 Bad Request, got %d", resp.StatusCode)
		}

		if ct := resp.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
			t.Errorf("Expected Content-Type 'application/json; charset=utf-8', got '%s'", ct)
		}

		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		if _, ok := errResp["error"]; !ok {
			t.Errorf("Expected error message in JSON, got: %+v", errResp)
		}
	})
}

func TestTransactionsAPI(t *testing.T) {
	server, ledgerService := setupTestServer()
	defer server.Close()

	t.Cleanup(func() {
		ledgerService.Reset()
	})

	budgetPayload := `{"category":"еда","limit":500}`
	http.Post(server.URL+"/api/budgets", "application/json", bytes.NewBufferString(budgetPayload))

	t.Run("Create transaction success", func(t *testing.T) {
		txPayload := `{"amount":450,"category":"еда","description":"ланч","date":"2025-09-10"}`
		resp, err := http.Post(server.URL+"/api/transactions", "application/json", bytes.NewBufferString(txPayload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201 Created, got %d", resp.StatusCode)
		}

		resp, err = http.Get(server.URL + "/api/transactions")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()
		var txs []TransactionResponse
		json.NewDecoder(resp.Body).Decode(&txs)
		if len(txs) != 1 || txs[0].Amount != 450 {
			t.Errorf("Transaction not found in list: %+v", txs)
		}
	})

	t.Run("Budget exceeded", func(t *testing.T) {
		txPayload := `{"amount":100,"category":"еда","description":"кофе","date":"2025-09-11"}`
		resp, err := http.Post(server.URL+"/api/transactions", "application/json", bytes.NewBufferString(txPayload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409 Conflict, got %d", resp.StatusCode)
		}
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errResp["error"] != "budget exceeded" {
			t.Errorf("Expected error 'budget exceeded', got '%s'", errResp["error"])
		}
	})

	t.Run("Bad JSON", func(t *testing.T) {
		txPayload := `{"amount":100, "category":"еда"`
		resp, err := http.Post(server.URL+"/api/transactions", "application/json", bytes.NewBufferString(txPayload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 Bad Request, got %d", resp.StatusCode)
		}
	})
}
