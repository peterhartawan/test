package entity_test

import (
	"amartha-reconciliation/internal/entity"
	"testing"
	"time"
)

func TestReconcileResult_ShowResult(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		fields entity.ReconcileResult
	}{
		{
			name: "valid result",
			fields: entity.ReconcileResult{
				TotalProcessed:   100,
				TotalMatch:       50,
				TotalUnmatched:   20,
				TotalDiscrepancy: 100.2500,
				MissingInBank: []entity.SystemTransaction{
					{
						TxID:            "123456",
						Amount:          100.00,
						Type:            entity.TYPE_DEBIT,
						TransactionTime: time.Time{},
					},
				},
				MissingInSystem: map[string][]entity.BankTransaction{
					"bank1.csv": {
						{
							UniqueIdentifier: "123456",
							Amount:           100.00,
							Date:             time.Time{},
						},
					},
				},
			},
		},
		{
			name: "valid result with missing in bank or system file none",
			fields: entity.ReconcileResult{
				TotalProcessed:   100,
				TotalMatch:       50,
				TotalUnmatched:   20,
				TotalDiscrepancy: 100.2500,
				MissingInBank:    []entity.SystemTransaction{},
				MissingInSystem:  map[string][]entity.BankTransaction{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			rr := tt.fields
			rr.ShowResult()
		})
	}
}
