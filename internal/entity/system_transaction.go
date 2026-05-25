package entity

import (
	"time"
)

const (
	TYPE_DEBIT  = "DEBIT"
	TYPE_CREDIT = "CREDIT"

	SYSTEM_FILE_HEADER = "trxID,amount,type,transactionTime"
)

const (
	SYSTEM_FILE_HEADER_TRX_ID = iota
	SYSTEM_FILE_HEADER_AMOUNT
	SYSTEM_FILE_HEADER_TYPE
	SYSTEM_FILE_HEADER_DATE
)

type SystemTransaction struct {
	TxID            string
	Amount          float64
	Type            string
	TransactionTime time.Time
}

func (st *SystemTransaction) GetAmountWithSign() float64 {
	if st.Type == TYPE_DEBIT {
		return -st.Amount
	}
	return st.Amount
}
