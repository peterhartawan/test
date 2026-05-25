package reconciler_test

import (
	"amartha-reconciliation/internal/entity"
	"amartha-reconciliation/internal/reconciler"
	"reflect"
	"testing"
	"time"
)

func TestParseSystemTransactions(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		systemFile string
		startDate  time.Time
		endDate    time.Time
		want       map[string]entity.SystemTransaction
		wantErr    bool
	}{
		{
			name:       "invalid system file failed open",
			systemFile: "./../../testdata/system/not_found.csv",
			startDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			want:       map[string]entity.SystemTransaction{},
			wantErr:    true,
		},
		{
			name:       "invalid system file row column",
			systemFile: "./../../testdata/system/malformed_row.csv",
			startDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			want:       map[string]entity.SystemTransaction{},
			wantErr:    true,
		},
		{
			name:       "invalid system file header",
			systemFile: "./../../testdata/system/malformed_header.csv",
			startDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			want:       map[string]entity.SystemTransaction{},
			wantErr:    true,
		},
		{
			name:       "invalid system file row amount",
			systemFile: "./../../testdata/system/malformed_amount.csv",
			startDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			want:       map[string]entity.SystemTransaction{},
			wantErr:    true,
		},
		{
			name:       "invalid system file row type",
			systemFile: "./../../testdata/system/malformed_type.csv",
			startDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			want:       map[string]entity.SystemTransaction{},
			wantErr:    true,
		},
		{
			name:       "invalid system file row transaction time",
			systemFile: "./../../testdata/system/malformed_transaction_time.csv",
			startDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			want:       map[string]entity.SystemTransaction{},
			wantErr:    true,
		},
		{
			name:       "invalid system file row duplicate trxID",
			systemFile: "./../../testdata/system/malformed_duplicate_trx.csv",
			startDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			want:       map[string]entity.SystemTransaction{},
			wantErr:    true,
		},
		{
			name:       "valid system file",
			systemFile: "./../../testdata/system/valid.csv",
			startDate:  time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			want: map[string]entity.SystemTransaction{
				"BCA002": {
					TxID:            "BCA002",
					Amount:          250000.75,
					Type:            entity.TYPE_DEBIT,
					TransactionTime: time.Date(2026, 1, 7, 11, 0, 0, 0, time.UTC),
				},
				"MISSING001": {
					TxID:            "MISSING001",
					Amount:          150000,
					Type:            entity.TYPE_DEBIT,
					TransactionTime: time.Date(2026, 1, 7, 11, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := reconciler.ParseSystemTransactions(tt.systemFile, tt.startDate, tt.endDate)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseSystemTransactions() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParseSystemTransactions() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSystemTransactions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckBankTransactions(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		bankFile           string
		startDate          time.Time
		endDate            time.Time
		systemTransactions map[string]entity.SystemTransaction
		matchResults       map[string]entity.MatchResult
		want               map[string]entity.MatchResult
		wantErr            bool
	}{
		{
			name:               "invalid bank file failed open",
			bankFile:           "./../../testdata/bank/not_found.csv",
			startDate:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:            time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			systemTransactions: map[string]entity.SystemTransaction{},
			matchResults:       map[string]entity.MatchResult{},
			want:               map[string]entity.MatchResult{},
			wantErr:            true,
		},
		{
			name:               "invalid bank file row column",
			bankFile:           "./../../testdata/bank/malformed_row.csv",
			startDate:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:            time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			systemTransactions: map[string]entity.SystemTransaction{},
			matchResults:       map[string]entity.MatchResult{},
			want:               map[string]entity.MatchResult{},
			wantErr:            true,
		},
		{
			name:               "invalid bank file header",
			bankFile:           "./../../testdata/bank/malformed_header.csv",
			startDate:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:            time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			systemTransactions: map[string]entity.SystemTransaction{},
			matchResults:       map[string]entity.MatchResult{},
			want:               map[string]entity.MatchResult{},
			wantErr:            true,
		},
		{
			name:               "invalid bank file amount",
			bankFile:           "./../../testdata/bank/malformed_amount.csv",
			startDate:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:            time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			systemTransactions: map[string]entity.SystemTransaction{},
			matchResults:       map[string]entity.MatchResult{},
			want:               map[string]entity.MatchResult{},
			wantErr:            true,
		},
		{
			name:               "invalid bank file date",
			bankFile:           "./../../testdata/bank/malformed_date.csv",
			startDate:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:            time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			systemTransactions: map[string]entity.SystemTransaction{},
			matchResults:       map[string]entity.MatchResult{},
			want:               map[string]entity.MatchResult{},
			wantErr:            true,
		},
		{
			name:      "invalid bank file duplicate unique_identifier",
			bankFile:  "./../../testdata/bank/malformed_duplicate_unique_identifier.csv",
			startDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			systemTransactions: map[string]entity.SystemTransaction{
				"BCA001": {
					TxID:            "BCA001",
					Amount:          100000.50,
					Type:            entity.TYPE_CREDIT,
					TransactionTime: time.Date(2026, 1, 1, 11, 0, 0, 0, time.UTC),
				},
			},
			matchResults: map[string]entity.MatchResult{},
			want:         map[string]entity.MatchResult{},
			wantErr:      true,
		},
		{
			name:      "valid",
			bankFile:  "./../../testdata/bank/valid.csv",
			startDate: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			systemTransactions: map[string]entity.SystemTransaction{
				"BCA002": {
					TxID:            "BCA002",
					Amount:          90000.50,
					Type:            entity.TYPE_CREDIT,
					TransactionTime: time.Date(2026, 1, 2, 11, 0, 0, 0, time.UTC),
				},
				"BCA003": {
					TxID:            "BCA003",
					Amount:          250000.75,
					Type:            entity.TYPE_CREDIT,
					TransactionTime: time.Date(2026, 1, 7, 11, 0, 0, 0, time.UTC),
				},
			},
			matchResults: map[string]entity.MatchResult{
				"BNI001": {
					Type:       entity.MATCH_RESULT_TYPE_MATCH,
					SourceFile: "./../../testdata/bank/valid-other-bank.csv",
					SystemTransaction: entity.SystemTransaction{
						TxID:            "BNI001",
						Amount:          10000.00,
						Type:            entity.TYPE_CREDIT,
						TransactionTime: time.Date(2026, 1, 2, 11, 0, 0, 0, time.UTC),
					},
					BankTransaction: entity.BankTransaction{
						UniqueIdentifier: "BNI001",
						Amount:           10000.00,
						Date:             time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
					},
					Discrepancy: 0,
				},
			},
			want: map[string]entity.MatchResult{
				"BNI001": {
					Type:       entity.MATCH_RESULT_TYPE_MATCH,
					SourceFile: "./../../testdata/bank/valid-other-bank.csv",
					SystemTransaction: entity.SystemTransaction{
						TxID:            "BNI001",
						Amount:          10000.00,
						Type:            entity.TYPE_CREDIT,
						TransactionTime: time.Date(2026, 1, 2, 11, 0, 0, 0, time.UTC),
					},
					BankTransaction: entity.BankTransaction{
						UniqueIdentifier: "BNI001",
						Amount:           10000.00,
						Date:             time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
					},
					Discrepancy: 0,
				},
				"BCA002": {
					Type:       entity.MATCH_RESULT_TYPE_MATCH,
					SourceFile: "./../../testdata/bank/valid.csv",
					SystemTransaction: entity.SystemTransaction{
						TxID:            "BCA002",
						Amount:          90000.50,
						Type:            entity.TYPE_CREDIT,
						TransactionTime: time.Date(2026, 1, 2, 11, 0, 0, 0, time.UTC),
					},
					BankTransaction: entity.BankTransaction{
						UniqueIdentifier: "BCA002",
						Amount:           90000.50,
						Date:             time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
					},
					Discrepancy: 0,
				},
				"BCA-X-001": {
					Type:              entity.MATCH_RESULT_TYPE_MISSING_IN_SYSTEM,
					SourceFile:        "./../../testdata/bank/valid.csv",
					SystemTransaction: entity.SystemTransaction{},
					BankTransaction: entity.BankTransaction{
						UniqueIdentifier: "BCA-X-001",
						Amount:           999999.99,
						Date:             time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC),
					},
					Discrepancy: 0,
				},
				"BCA003": {
					Type:       entity.MATCH_RESULT_TYPE_UNMATCH,
					SourceFile: "./../../testdata/bank/valid.csv",
					SystemTransaction: entity.SystemTransaction{
						TxID:            "BCA003",
						Amount:          250000.75,
						Type:            entity.TYPE_CREDIT,
						TransactionTime: time.Date(2026, 1, 7, 11, 0, 0, 0, time.UTC),
					},
					BankTransaction: entity.BankTransaction{
						UniqueIdentifier: "BCA003",
						Amount:           -250000.75,
						Date:             time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
					},
					Discrepancy: 0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := reconciler.CheckBankTransactions(tt.bankFile, tt.startDate, tt.endDate, tt.systemTransactions, tt.matchResults)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CheckBankTransactions() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CheckBankTransactions() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(tt.matchResults, tt.want) {
				t.Errorf("CheckBankTransactions() = %v, want %v", tt.matchResults, tt.want)
			}
		})
	}
}

func TestReconcile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		input   entity.InputData
		want    entity.ReconcileResult
		wantErr bool
	}{
		{
			name: "invalid parse system transactions failed",
			input: entity.InputData{
				SystemFile: "./../../testdata/system/malformed_row.csv",
				BankFiles:  []string{"./../../testdata/bank/valid.csv"},
				StartDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:    time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			},
			want:    entity.ReconcileResult{},
			wantErr: true,
		},
		{
			name: "invalid check bank transactions failed",
			input: entity.InputData{
				SystemFile: "./../../testdata/system/valid.csv",
				BankFiles:  []string{"./../../testdata/bank/malformed_row.csv"},
				StartDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:    time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			},
			want:    entity.ReconcileResult{},
			wantErr: true,
		},
		{
			name: "valid",
			input: entity.InputData{
				SystemFile: "./../../testdata/system/valid.csv",
				BankFiles:  []string{"./../../testdata/bank/valid.csv"},
				StartDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:    time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
			},
			want: entity.ReconcileResult{
				TotalProcessed:   3,
				TotalMatch:       1,
				TotalUnmatched:   1,
				TotalDiscrepancy: 160000.25,
				MissingInBank: []entity.SystemTransaction{
					{
						TxID:            "MISSING001",
						Amount:          150000,
						Type:            entity.TYPE_DEBIT,
						TransactionTime: time.Date(2026, 1, 7, 11, 0, 0, 0, time.UTC),
					},
				},
				MissingInSystem: map[string][]entity.BankTransaction{
					"./../../testdata/bank/valid.csv": {
						{
							UniqueIdentifier: "BCA-X-001",
							Amount:           999999.99,
							Date:             time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC),
						},
						{
							UniqueIdentifier: "BCA003",
							Amount:           -250000.75,
							Date:             time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := reconciler.Reconcile(tt.input)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Reconcile() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Reconcile() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reconcile() = %v, want %v", got, tt.want)
			}
		})
	}
}
