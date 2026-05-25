package main

import (
	"flag"
	"log"

	"amartha-reconciliation/internal/entity"
	"amartha-reconciliation/internal/reconciler"
)

func main() {
	systemFile := flag.String("system", "", "system transaction csv file")
	bankFiles := flag.String("bank", "", "bank transaction csv file (separated by comma)")
	startDate := flag.String("start", "", "start date")
	endDate := flag.String("end", "", "end date")
	flag.Parse()

	inputRequest := entity.InputRequest{
		SystemFile: *systemFile,
		BankFiles:  *bankFiles,
		StartDate:  *startDate,
		EndDate:    *endDate,
	}

	inputData, err := inputRequest.ParseAndValidate()
	if err != nil {
		log.Fatalf("input error: %v", err)
		return
	}

	reconcileResult, err := reconciler.Reconcile(inputData)
	if err != nil {
		log.Fatalf("reconciliation failed: %v", err)
		return
	}

	reconcileResult.ShowResult()
}
