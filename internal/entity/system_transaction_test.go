package entity_test

import (
	"amartha-reconciliation/internal/entity"
	"testing"
	"time"
)

func TestSystemTransaction_GetAmountWithSign(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		st   entity.SystemTransaction
		want float64
	}{
		{
			name: "valid debit type",
			st: entity.SystemTransaction{
				TxID:            "BCA002",
				Amount:          250000.75,
				Type:            entity.TYPE_DEBIT,
				TransactionTime: time.Date(2026, 1, 7, 11, 0, 0, 0, time.UTC),
			},
			want: -250000.75,
		},
		{
			name: "valid credit type",
			st: entity.SystemTransaction{
				TxID:            "BCA003",
				Amount:          100000.00,
				Type:            entity.TYPE_CREDIT,
				TransactionTime: time.Date(2026, 1, 7, 11, 0, 0, 0, time.UTC),
			},
			want: 100000.00,
		},
		{
			name: "valid credit type anomaly negative sign",
			st: entity.SystemTransaction{
				TxID:            "BCA003",
				Amount:          -100000.00,
				Type:            entity.TYPE_CREDIT,
				TransactionTime: time.Date(2026, 1, 7, 11, 0, 0, 0, time.UTC),
			},
			want: -100000.00,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			st := tt.st
			got := st.GetAmountWithSign()
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("GetAmountWithSign() = %v, want %v", got, tt.want)
			}
		})
	}
}
