package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

type Validatable interface {
	Validate() error
}

// Transaction - описывает одну финансовую операцию.
type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        time.Time
}

// Budget - описывает лимит трат по категории.
type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
}

// Ledger - наш главный сервис, управляющий финансами.
type Ledger struct {
	transactions []Transaction
	budgets      map[string]Budget
	nextTxID     int
}

// --- Методы Ledger ---

// NewLedger - конструктор для Ledger.
func NewLedger() *Ledger {
	return &Ledger{
		transactions: make([]Transaction, 0),
		budgets:      make(map[string]Budget),
		nextTxID:     1, // Начинаем ID с 1
	}
}

// AddTransaction - добавляет транзакцию в Ledger.
func (l *Ledger) AddTransaction(tx Transaction) error {
	if err := tx.Validate(); err != nil {
		return fmt.Errorf("невалидная транзакция: %w", err)
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
			return fmt.Errorf("превышен бюджет для категории '%s'", tx.Category)
		}
	}

	tx.ID = l.nextTxID
	l.nextTxID++
	tx.Date = time.Now()

	l.transactions = append(l.transactions, tx)
	return nil
}

// SetBudget - устанавливает бюджет для категории.
func (l *Ledger) SetBudget(b Budget) error {
	if err := b.Validate(); err != nil {
		return fmt.Errorf("невалидный бюджет: %w", err)
	}
	l.budgets[b.Category] = b
	return nil
}

// LoadBudgets - загружает бюджеты из JSON.
func (l *Ledger) LoadBudgets(r io.Reader) error {
	var budgets []Budget
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&budgets); err != nil {
		return fmt.Errorf("ошибка декодирования JSON: %w", err)
	}
	for _, b := range budgets {
		l.SetBudget(b)
	}
	return nil
}

// Validate - реализует интерфейс Validatable для Transaction.
func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return fmt.Errorf("сумма транзакции должна быть положительной, получено: %.2f", t.Amount)
	}
	if t.Category == "" {
		return errors.New("категория транзакции не может быть пустой")
	}
	return nil
}

// Validate - реализует интерфейс Validatable для Budget.
func (b *Budget) Validate() error {
	if b.Limit <= 0 {
		return fmt.Errorf("лимит бюджета должен быть положительным, получено: %.2f", b.Limit)
	}
	if b.Category == "" {
		return errors.New("категория бюджета не может быть пустой")
	}
	return nil // Ошибки нет
}

// CheckValid - демонстрирует полиморфизм через интерфейс Validatable.
func CheckValid(v Validatable) {
	fmt.Printf("--> Проверяем %+v\n", v)
	if err := v.Validate(); err != nil {
		fmt.Printf("    Результат: НЕВАЛИДНО! Ошибка: %v\n", err)
	} else {
		fmt.Println("    Результат: ВАЛИДНО!")
	}
}

func main() {
	fmt.Println("--- Демонстрация работы интерфейса Validatable ---")
	validTx := Transaction{Amount: 100, Category: "Еда"}
	invalidTxAmount := Transaction{Amount: -50, Category: "Еда"}
	invalidTxCategory := Transaction{Amount: 100, Category: ""}

	validBudget := Budget{Limit: 5000, Category: "Развлечения"}
	invalidBudgetLimit := Budget{Limit: 0, Category: "Развлечения"}
	invalidBudgetCategory := Budget{Limit: 5000, Category: ""}

	// Передаем их все в одну и ту же функцию
	CheckValid(&validTx)
	CheckValid(&invalidTxAmount)
	CheckValid(&invalidTxCategory)
	fmt.Println()
	CheckValid(&validBudget)
	CheckValid(&invalidBudgetLimit)
	CheckValid(&invalidBudgetCategory)
	fmt.Println("---------------------------------------------------")

	fmt.Println("--- Тестирование основной логики Ledger ---")
	ledger := NewLedger()

	// Попытка установить невалидный бюджет напрямую
	err := ledger.SetBudget(Budget{Category: "Путешествия", Limit: -100})
	if err != nil {
		fmt.Printf("Ожидаемая ошибка при установке бюджета: %v\n", err)
	}
	fmt.Println("--------------------------------")

	// Попытка добавить невалидную транзакцию
	err = ledger.AddTransaction(Transaction{Category: "", Amount: 200})
	if err != nil {
		fmt.Printf("Ожидаемая ошибка при добавлении транзакции: %v\n", err)
	}
	fmt.Println("--------------------------------")

	// Успешный сценарий
	err = ledger.SetBudget(Budget{Category: "Еда", Limit: 500})
	if err != nil {
		fmt.Printf("Неожиданная ошибка: %v\n", err)
	} else {
		fmt.Println("Бюджет 'Еда' успешно установлен.")
	}

	err = ledger.AddTransaction(Transaction{Category: "Еда", Amount: 300, Description: "Обед"})
	if err != nil {
		fmt.Printf("Неожиданная ошибка: %v\n", err)
	} else {
		fmt.Println("Транзакция 'Обед' успешно добавлена.")
	}

	// Сценарий превышения бюджета
	err = ledger.AddTransaction(Transaction{Category: "Еда", Amount: 250, Description: "Ужин"})
	if err != nil {
		fmt.Printf("Ожидаемая ошибка превышения бюджета: %v\n", err)
	}
	fmt.Println("--------------------------------")

	fmt.Println("\nИтоговое состояние Ledger:")
	fmt.Printf("Транзакции: %+v\n", ledger.transactions)
	fmt.Printf("Бюджеты: %+v\n", ledger.budgets)
}
