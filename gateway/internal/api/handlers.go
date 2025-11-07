package api

import (
	"encoding/json"
	"errors"
	"ledger"
	"log"
	"net/http"
	"time"
)

type Handlers struct {
	ledgerService *ledger.Ledger
}

func NewHandlers(ledgerService *ledger.Ledger) *Handlers {
	return &Handlers{ledgerService: ledgerService}
}

// --- Вспомогательные функции для ответов ---

// respondWithError - отправляет JSON-ответ с ошибкой и заданным HTTP-статусом.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON - отправляет JSON-ответ с данными и статусом.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

// --- Middleware ---

// LoggingMiddleware - логирует информацию о каждом запросе.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"%s %s %v",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

// --- Обработчики Транзакций ---

// CreateTransaction - обработчик для POST /api/transactions
func (h *Handlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Преобразуем DTO в доменную модель.
	tx, err := ToDomainTransaction(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
		return
	}

	// Вызываем метод доменного сервиса.
	createdTx, err := h.ledgerService.AddTransaction(tx)
	if err != nil {
		if errors.Is(err, ledger.ErrBudgetExceeded) {
			respondWithError(w, http.StatusConflict, "budget exceeded")
		} else if err.Error() == "invalid transaction" {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			log.Printf("Internal server error: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, ToTransactionResponse(createdTx))
}

// ListTransactions - обработчик для GET /api/transactions
func (h *Handlers) ListTransactions(w http.ResponseWriter, r *http.Request) {
	transactions := h.ledgerService.ListTransactions()

	responses := make([]TransactionResponse, 0, len(transactions))
	for _, tx := range transactions {
		responses = append(responses, ToTransactionResponse(tx))
	}

	respondWithJSON(w, http.StatusOK, responses)
}

// --- Обработчики Бюджетов ---

// CreateOrUpdateBudget - обработчик для POST /api/budgets
func (h *Handlers) CreateOrUpdateBudget(w http.ResponseWriter, r *http.Request) {
	var req CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	budget := ToDomainBudget(req)

	if err := h.ledgerService.SetBudget(budget); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, ToBudgetResponse(budget))
}

// ListBudgets - обработчик для GET /api/budgets
func (h *Handlers) ListBudgets(w http.ResponseWriter, r *http.Request) {
	budgets := h.ledgerService.ListBudgets()

	responses := make([]BudgetResponse, 0, len(budgets))
	for _, b := range budgets {
		responses = append(responses, ToBudgetResponse(b))
	}

	respondWithJSON(w, http.StatusOK, responses)
}
