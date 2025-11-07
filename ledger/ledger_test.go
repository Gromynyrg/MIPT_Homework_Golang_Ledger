package ledger

import (
	"errors"
	"testing"
)

// TestTransaction_Validate - тестирует метод Validate() у Transaction.
func TestTransaction_Validate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		tx      Transaction
		wantErr bool
	}{
		{
			name:    "valid transaction",
			tx:      Transaction{Amount: 100, Category: "food"},
			wantErr: false,
		},
		{
			name:    "invalid zero amount",
			tx:      Transaction{Amount: 0, Category: "food"},
			wantErr: true,
		},
		{
			name:    "invalid negative amount",
			tx:      Transaction{Amount: -50, Category: "food"},
			wantErr: true,
		},
		{
			name:    "invalid empty category",
			tx:      Transaction{Amount: 100, Category: ""},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.tx.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Transaction.Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// TestBudget_Validate - тестирует метод Validate() у Budget.
func TestBudget_Validate(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		budget  Budget
		wantErr bool
	}{
		{name: "valid budget", budget: Budget{Category: "travel", Limit: 1000}, wantErr: false},
		{name: "invalid zero limit", budget: Budget{Category: "travel", Limit: 0}, wantErr: true},
		{name: "invalid negative limit", budget: Budget{Category: "travel", Limit: -100}, wantErr: true},
		{name: "invalid empty category", budget: Budget{Category: "", Limit: 1000}, wantErr: true},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.budget.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Budget.Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// TestLedger_AddTransaction_BudgetExceeded - тестирует бизнес-правило превышения бюджета.
func TestLedger_AddTransaction_BudgetExceeded(t *testing.T) {

	ledger := NewLedger()
	t.Cleanup(func() {
		ledger.Reset()
	})

	err := ledger.SetBudget(Budget{Category: "еда", Limit: 500})
	if err != nil {
		t.Fatalf("SetBudget() failed: %v", err)
	}

	// 2. Выполнение и проверка (Act & Assert)

	_, err = ledger.AddTransaction(Transaction{Amount: 300, Category: "еда"})
	if err != nil {
		t.Fatalf("AddTransaction() unexpected error for valid tx: %v", err)
	}

	if len(ledger.ListTransactions()) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(ledger.ListTransactions()))
	}

	_, err = ledger.AddTransaction(Transaction{Amount: 250, Category: "еда"})

	if !errors.Is(err, ErrBudgetExceeded) {
		t.Errorf("Expected error %v, got %v", ErrBudgetExceeded, err)
	}

	if len(ledger.ListTransactions()) != 1 {
		t.Errorf("Transaction that exceeds budget should not be added. Expected 1 transaction, got %d", len(ledger.ListTransactions()))
	}
}
