package api

import (
	"ledger"
	"time"
)

// --- Transaction DTOs ---

// CreateTransactionRequest - структура для парсинга JSON-тела запроса на создание транзакции.
type CreateTransactionRequest struct {
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

// TransactionResponse - структура для отправки данных о транзакции клиенту.
type TransactionResponse struct {
	ID          int       `json:"id"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

// --- Budget DTOs ---

// CreateBudgetRequest - DTO для создания/обновления бюджета.
type CreateBudgetRequest struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
}

// BudgetResponse - DTO для ответа с данными о бюджете.
type BudgetResponse struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
}

// --- Функции-преобразователи (Mappers) ---

// ToDomainTransaction - преобразует DTO в доменную модель Transaction.
func ToDomainTransaction(req CreateTransactionRequest) (ledger.Transaction, error) {
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return ledger.Transaction{}, err
	}

	return ledger.Transaction{
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
		Date:        parsedDate,
	}, nil
}

// ToTransactionResponse - преобразует доменную модель в DTO для ответа.
func ToTransactionResponse(tx ledger.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:          tx.ID,
		Amount:      tx.Amount,
		Category:    tx.Category,
		Description: tx.Description,
		Date:        tx.Date,
	}
}

// ToDomainBudget - преобразует DTO в доменную модель Budget.
func ToDomainBudget(req CreateBudgetRequest) ledger.Budget {
	return ledger.Budget{
		Category: req.Category,
		Limit:    req.Limit,
	}
}

// ToBudgetResponse - преобразует доменную модель в DTO для ответа.
func ToBudgetResponse(b ledger.Budget) BudgetResponse {
	return BudgetResponse{
		Category: b.Category,
		Limit:    b.Limit,
	}
}
