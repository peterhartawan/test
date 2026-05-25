# Amartha Reconciliation Tool

A simple CLI-based reconciliation tool built with Go.

This tool compares internal system transactions against bank transaction files and generates reconciliation results including:
- matched transactions
- unmatched transactions
- discrepancy totals
- missing transactions in system
- missing transactions in bank files

---

# Features

- Stream-based CSV processing 
- Supports multiple bank files
- Duplicate transaction detection
- Date-range filtering
- Transaction reconciliation
- Discrepancy calculation

---

# Assumptions

## 1. Transaction Identifier Uniqueness

`trxID` and `unique_identifier` are assumed to be globally unique across all provided bank files within a reconciliation run.

Duplicate transaction identifiers will return an error.

---

## 2. System Transaction Amount

System transaction amounts are always stored as positive values.

Transaction direction is determined using the `type` field:
- `CREDIT`
- `DEBIT`

---

## 3. Bank Transaction Amount

Bank transaction amounts may contain negative values.

Negative values are interpreted as transaction direction indicators.

Example:

```txt
-250000
```

means:
- money out
- debit transaction

---

## 4. Reconciliation Logic

### Matched Transaction

Transaction exists in both:
- system file
- bank file

AND:
- amount matches
- transaction direction/sign is valid

---

### Unmatched Transaction

Transaction exists in both files but reconciliation fails semantically.

Examples:
- amount mismatch
- sign mismatch
- direction mismatch

---

### Discrepancy

Discrepancy represents the monetary difference only.

Formula:

```txt
abs(systemAmount - abs(bankAmount))
```

Examples:

```txt
System: 250000
Bank:   240000

Discrepancy = 10000
```

```txt
System: 250000
Bank:  -250000

Discrepancy = 0
```

The second example is considered unmatched due to sign mismatch, but monetary discrepancy remains zero because the amount magnitude is identical.

---

## 5. Date Filtering

Date filtering uses:
- inclusive start date
- inclusive end date

Internally, the implementation uses:
- inclusive lower bound
- exclusive upper bound

Example:

```txt
2026-01-01 <= transactionTime < 2026-01-08
```

for:
```txt
--start=2026-01-01
--end=2026-01-07
```

---

## 6. CSV Header Validation

CSV headers must match expected formats exactly.

### System Transaction CSV

```csv
trxID,amount,type,transactionTime
```

### Bank Transaction CSV

```csv
unique_identifier,amount,date
```

---

# Usage

## Run

```bash
go run ./cmd/tools/main.go \
  --system=./input/system_transaction.csv \
  --bank=./input/bank_bca.csv,./input/bank_bni.csv \
  --start=2026-01-01 \
  --end=2026-01-07
```

---

# Example Output

```txt
======================
RECONCILIATION SUMMARY
======================

Processed Transactions : 16
Matched Transactions   : 9
Unmatched Transactions : 3
Total Discrepancy      : 55000.25

Missing In System File:
- BCA-X-001
- BCA-X-002
- JAGO-X-001
- BNI-X-001
- BNI-X-002

Missing In Bank File:
- ROUND003
- ROUND002
- MISSING001
- ROUND001
```

---

# Unit Test

Run all tests:

```bash
go test ./...
```

Run with coverage:

```bash
go test ./... -cover
```

---
