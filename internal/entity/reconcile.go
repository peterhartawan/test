package entity

import "fmt"

const (
	MATCH_RESULT_TYPE_MATCH = iota
	MATCH_RESULT_TYPE_UNMATCH
	MATCH_RESULT_TYPE_MISSING_IN_SYSTEM
	MATCH_RESULT_TYPE_MISSING_IN_BANK
)

type MatchResult struct {
	Type              int
	SourceFile        string
	SystemTransaction SystemTransaction
	BankTransaction   BankTransaction
	Discrepancy       float64
}

type ReconcileResult struct {
	TotalProcessed   int
	TotalMatch       int
	TotalUnmatched   int
	TotalDiscrepancy float64
	MissingInBank    []SystemTransaction          // Exists at system transaction only
	MissingInSystem  map[string][]BankTransaction // Exists at bank transaction only
}

func (rr *ReconcileResult) ShowResult() {
	fmt.Println("======================")
	fmt.Println("RECONCILIATION SUMMARY")
	fmt.Println("======================")
	fmt.Println()

	fmt.Printf("Processed Transactions : %d\n", rr.TotalProcessed)
	fmt.Printf("Matched Transactions   : %d\n", rr.TotalMatch)
	fmt.Printf("Unmatched Transactions : %d\n", rr.TotalUnmatched)
	fmt.Printf("Total Discrepancy      : %.2f\n", rr.TotalDiscrepancy)

	fmt.Println()

	fmt.Println("Missing In System File:")
	if len(rr.MissingInSystem) == 0 {
		fmt.Println("- None")
	} else {
		for _, bankTransactions := range rr.MissingInSystem {
			for _, bankTransaction := range bankTransactions {
				fmt.Printf("- %s\n", bankTransaction.UniqueIdentifier)
			}
		}
	}

	fmt.Println()

	fmt.Println("Missing In Bank File:")
	if len(rr.MissingInBank) == 0 {
		fmt.Println("- None")
	} else {
		for _, systemTransaction := range rr.MissingInBank {
			fmt.Printf("- %s\n", systemTransaction.TxID)
		}
	}
}
