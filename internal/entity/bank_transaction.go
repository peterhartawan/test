package entity

import "time"

const (
	BANK_FILE_HEADER = "unique_identifier,amount,date"
)

const (
	BANK_FILE_HEADER_UNIQUE_IDENTIFIER = iota
	BANK_FILE_HEADER_AMOUNT
	BANK_FILE_HEADER_DATE
)

type BankTransaction struct {
	UniqueIdentifier string
	Amount           float64
	Date             time.Time
}
