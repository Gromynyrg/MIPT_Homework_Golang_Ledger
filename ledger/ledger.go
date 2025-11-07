package ledger

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Validatable - это интерфейс для структур, которые можно проверить на корректность.
type Validatable interface {
	Validate() error
}

// Transaction - описывает одну финансовую операцию.
// Поля сделаны экспортируемыми (начинаются с большой буквы), чтобы к ним можно было
// обращаться из других пакетов, например, из нашего gateway.
type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        time.Time
}

// Budget - описывает лимит трат по категории.
type Budget struct {
	Category string
	Limit    float64
}

// Ledger - наш главный сервис, управляющий финансами.
// Мы добавили мьютекс (sync.Mutex) для безопасной работы с данными
// из нескольких горутин (что происходит при обработке нескольких HTTP-запросов одновременно).
type Ledger struct {
	transactions []Transaction
	budgets      map[string]Budget
	nextTxID     int
	mu           sync.Mutex // Мьютекс для защиты доступа к срезу и карте
}

// NewLedger - конструктор для Ledger.
func NewLedger() *Ledger {
	return &Ledger{
		transactions: make([]Transaction, 0),
		budgets:      make(map[string]Budget),
		nextTxID:     1, // Начинаем ID с 1
	}
}

// Определим кастомные ошибки для более точной обработки в HTTP-слое.
var ErrBudgetExceeded = errors.New("budget exceeded")

// AddTransaction - добавляет транзакцию в Ledger.
func (l *Ledger) AddTransaction(tx Transaction) (Transaction, error) {
	// Блокируем доступ к данным на время выполнения функции, чтобы избежать гонок данных.
	l.mu.Lock()
	// defer гарантирует, что мьютекс будет освобожден при выходе из функции,
	// даже если произойдет ошибка.
	defer l.mu.Unlock()

	if err := tx.Validate(); err != nil {
		return Transaction{}, fmt.Errorf("invalid transaction: %w", err)
	}

	budget, budgetExists := l.budgets[tx.Category]
	if budgetExists {
		currentSpend := 0.0
		for _, existingTx := range l.transactions {
			if existingTx.Category == tx.Category {
				currentSpend += existingTx.Amount
			}
		}
		if currentSpend+tx.Amount > budget.Limit {
			// Возвращаем нашу специальную ошибку.
			return Transaction{}, ErrBudgetExceeded
		}
	}

	tx.ID = l.nextTxID
	l.nextTxID++
	// В реальном приложении дату нужно брать из запроса, но по заданию пока так.
	if tx.Date.IsZero() {
		tx.Date = time.Now()
	}

	l.transactions = append(l.transactions, tx)
	// Возвращаем созданную транзакцию с присвоенным ID и датой.
	return tx, nil
}

// SetBudget - устанавливает бюджет для категории.
func (l *Ledger) SetBudget(b Budget) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if err := b.Validate(); err != nil {
		return fmt.Errorf("invalid budget: %w", err)
	}
	l.budgets[b.Category] = b
	return nil
}

// ListTransactions - возвращает все транзакции.
func (l *Ledger) ListTransactions() []Transaction {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Возвращаем копию среза, чтобы избежать изменений извне.
	result := make([]Transaction, len(l.transactions))
	copy(result, l.transactions)
	return result
}

// ListBudgets - возвращает все бюджеты.
func (l *Ledger) ListBudgets() []Budget {
	l.mu.Lock()
	defer l.mu.Unlock()

	result := make([]Budget, 0, len(l.budgets))
	for _, b := range l.budgets {
		result = append(result, b)
	}
	return result
}

// Validate - реализует интерфейс Validatable для Transaction.
func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return fmt.Errorf("transaction amount must be positive, got: %.2f", t.Amount)
	}
	if t.Category == "" {
		return errors.New("transaction category cannot be empty")
	}
	return nil
}

// Validate - реализует интерфейс Validatable для Budget.
func (b *Budget) Validate() error {
	if b.Limit <= 0 {
		return fmt.Errorf("budget limit must be positive, got: %.2f", b.Limit)
	}
	if b.Category == "" {
		return errors.New("budget category cannot be empty")
	}
	return nil
}

// Reset - сбрасывает состояние Ledger. ИСПОЛЬЗОВАТЬ ТОЛЬКО В ТЕСТАХ!
func (l *Ledger) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.transactions = make([]Transaction, 0)
	l.budgets = make(map[string]Budget)
	l.nextTxID = 1
}
