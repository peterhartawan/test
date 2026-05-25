package reconciler

import (
	"amartha-reconciliation/internal/entity"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseSystemTransactions(systemFile string, startDate, endDate time.Time) (map[string]entity.SystemTransaction, error) {
	systemTransactions := map[string]entity.SystemTransaction{}

	// Read system file content
	fileContent, err := os.Open(systemFile)
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	// Lazy read system file content
	csvReader := csv.NewReader(fileContent)
	row := 0
	for {
		row++
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("reading system file content failed, row %d, error: %w", row, err)
		}

		// Check header
		if row == 1 {
			headerStr := strings.Join(record, ",")
			if headerStr != entity.SYSTEM_FILE_HEADER {
				return nil, fmt.Errorf("system file header is invalid, file %s, got %s, expected %s", systemFile, headerStr, entity.SYSTEM_FILE_HEADER)
			}

			continue
		}

		// Parse amount
		rawAmount := record[entity.SYSTEM_FILE_HEADER_AMOUNT]
		amount, err := strconv.ParseFloat(rawAmount, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid amount at row %d, got %s, error: %w", row, rawAmount, err)
		}

		// Check type
		rawType := record[entity.SYSTEM_FILE_HEADER_TYPE]
		if rawType != entity.TYPE_DEBIT && rawType != entity.TYPE_CREDIT {
			return nil, fmt.Errorf("invalid type at row %d, got %s, expected %s or %s", row, rawType, entity.TYPE_DEBIT, entity.TYPE_CREDIT)
		}

		// Parse transaction time
		rawTransactionTime := record[entity.SYSTEM_FILE_HEADER_DATE]
		transactionTime, err := time.Parse("2006-01-02T15:04:05Z", rawTransactionTime)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction time at row %d, got %s, error: %w", row, rawTransactionTime, err)
		}

		// Skip transactions outside of the date range
		endDateBoundary := endDate.AddDate(0, 0, 1)
		if transactionTime.Before(startDate) || !transactionTime.Before(endDateBoundary) {
			continue
		}

		rawTxID := record[entity.SYSTEM_FILE_HEADER_TRX_ID]
		// Check duplicate data
		if _, ok := systemTransactions[rawTxID]; ok {
			return nil, fmt.Errorf("duplicate transaction id at row %d, got %s", row, rawTxID)
		}

		systemTransactions[rawTxID] = entity.SystemTransaction{
			TxID:            rawTxID,
			Amount:          amount,
			Type:            rawType,
			TransactionTime: transactionTime,
		}
	}

	return systemTransactions, nil
}

func CheckBankTransactions(bankFile string, startDate, endDate time.Time, systemTransactions map[string]entity.SystemTransaction, matchResults map[string]entity.MatchResult) error {
	// Read bank file content
	fileContent, err := os.Open(bankFile)
	if err != nil {
		return err
	}
	defer fileContent.Close()

	csvReader := csv.NewReader(fileContent)
	row := 0
	for {
		row++
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("reading bank file content failed, row %d, error: %w", row, err)
		}

		// Check header
		if row == 1 {
			headerStr := strings.Join(record, ",")
			if headerStr != entity.BANK_FILE_HEADER {
				return fmt.Errorf("bank file header is invalid, file %s, got %s, expected %s", bankFile, headerStr, entity.BANK_FILE_HEADER)
			}
			continue
		}

		// Parse amount
		rawAmount := record[entity.BANK_FILE_HEADER_AMOUNT]
		amount, err := strconv.ParseFloat(rawAmount, 64)
		if err != nil {
			return fmt.Errorf("invalid amount at row %d, got %s, error: %w", row, rawAmount, err)
		}

		// Parse transaction time
		rawDate := record[entity.BANK_FILE_HEADER_DATE]
		date, err := time.Parse("2006-01-02", rawDate)
		if err != nil {
			return fmt.Errorf("invalid transaction time at row %d, got %s, error: %w", row, rawDate, err)
		}

		// Skip transactions outside of the date range
		endDateBoundary := endDate.AddDate(0, 0, 1)
		if date.Before(startDate) || !date.Before(endDateBoundary) {
			continue
		}

		// Check transaction exist in system transactions
		rawUniqueIdentifier := record[entity.BANK_FILE_HEADER_UNIQUE_IDENTIFIER]
		systemTransaction, ok := systemTransactions[rawUniqueIdentifier]
		if !ok {
			matchResults[rawUniqueIdentifier] = entity.MatchResult{
				Type:              entity.MATCH_RESULT_TYPE_MISSING_IN_SYSTEM,
				SourceFile:        bankFile,
				SystemTransaction: entity.SystemTransaction{},
				BankTransaction: entity.BankTransaction{
					UniqueIdentifier: rawUniqueIdentifier,
					Amount:           amount,
					Date:             date,
				},
				Discrepancy: 0,
			}

			continue
		}

		matchType := entity.MATCH_RESULT_TYPE_MATCH
		// Compare amount
		if amount != systemTransaction.GetAmountWithSign() {
			matchType = entity.MATCH_RESULT_TYPE_UNMATCH
		}

		// Check duplicate data
		if _, ok := matchResults[rawUniqueIdentifier]; ok {
			return fmt.Errorf("duplicate unique identifier at row %d, got %s", row, rawUniqueIdentifier)
		}

		// Remove negative sign
		normalizedAmount := math.Abs(amount)
		matchResults[rawUniqueIdentifier] = entity.MatchResult{
			Type:              matchType,
			SourceFile:        bankFile,
			SystemTransaction: systemTransaction,
			BankTransaction: entity.BankTransaction{
				UniqueIdentifier: rawUniqueIdentifier,
				Amount:           amount,
				Date:             date,
			},
			Discrepancy: math.Abs(systemTransaction.Amount - normalizedAmount),
		}
	}

	return nil
}

func Reconcile(input entity.InputData) (entity.ReconcileResult, error) {
	systemTransactions, err := ParseSystemTransactions(input.SystemFile, input.StartDate, input.EndDate)
	if err != nil {
		return entity.ReconcileResult{}, err
	}

	matchedResult := map[string]entity.MatchResult{}
	for _, bankFile := range input.BankFiles {
		err := CheckBankTransactions(bankFile, input.StartDate, input.EndDate, systemTransactions, matchedResult)
		if err != nil {
			return entity.ReconcileResult{}, err
		}
	}

	for _, systemTransaction := range systemTransactions {
		if _, ok := matchedResult[systemTransaction.TxID]; !ok {
			matchedResult[systemTransaction.TxID] = entity.MatchResult{
				Type:              entity.MATCH_RESULT_TYPE_MISSING_IN_BANK,
				SourceFile:        input.SystemFile,
				SystemTransaction: systemTransaction,
				BankTransaction:   entity.BankTransaction{},
				Discrepancy:       0,
			}
		}
	}

	reconcileResult := entity.ReconcileResult{
		TotalProcessed:   len(systemTransactions),
		TotalMatch:       0,
		TotalUnmatched:   0,
		TotalDiscrepancy: 0,
		MissingInBank:    []entity.SystemTransaction{},
		MissingInSystem:  map[string][]entity.BankTransaction{},
	}
	for _, matchResult := range matchedResult {
		switch matchResult.Type {
		case entity.MATCH_RESULT_TYPE_MATCH:
			reconcileResult.TotalMatch++
		case entity.MATCH_RESULT_TYPE_UNMATCH:
			reconcileResult.TotalUnmatched++
		case entity.MATCH_RESULT_TYPE_MISSING_IN_BANK:
			reconcileResult.MissingInBank = append(reconcileResult.MissingInBank, matchResult.SystemTransaction)
		case entity.MATCH_RESULT_TYPE_MISSING_IN_SYSTEM:
			reconcileResult.MissingInSystem[matchResult.SourceFile] = append(reconcileResult.MissingInSystem[matchResult.SourceFile], matchResult.BankTransaction)
		}

		reconcileResult.TotalDiscrepancy += matchResult.Discrepancy
	}

	return reconcileResult, nil
}
