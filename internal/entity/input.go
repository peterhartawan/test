package entity

import (
	"errors"
	"os"
	"strings"
	"time"
)

type InputRequest struct {
	SystemFile string
	BankFiles  string
	StartDate  string
	EndDate    string
}

type InputData struct {
	SystemFile string
	BankFiles  []string
	StartDate  time.Time
	EndDate    time.Time
}

func (i *InputRequest) ParseAndValidate() (InputData, error) {
	// Validate system file
	if i.SystemFile == "" {
		return InputData{}, errors.New("system file is required")
	}
	if !strings.HasSuffix(i.SystemFile, ".csv") {
		return InputData{}, errors.New("system file must be csv format")
	}
	if _, err := os.Stat(i.SystemFile); err != nil {
		return InputData{}, errors.New("system file not found, file: " + i.SystemFile)
	}

	// Validate bank files
	if len(i.BankFiles) == 0 {
		return InputData{}, errors.New("bank file is required")
	}
	bankFilesList := strings.Split(i.BankFiles, ",")
	for _, bankFile := range bankFilesList {
		if !strings.HasSuffix(bankFile, ".csv") {
			return InputData{}, errors.New("bank file must be csv format, got " + bankFile)
		}
		if _, err := os.Stat(bankFile); err != nil {
			return InputData{}, errors.New("bank file not found, file: " + bankFile)
		}
	}

	// Validate start date
	if i.StartDate == "" {
		return InputData{}, errors.New("start date is required")
	}
	startDateTime, err := time.Parse("2006-01-02", i.StartDate)
	if err != nil {
		return InputData{}, errors.New("start date invalid format")
	}

	// Validate end date
	if i.EndDate == "" {
		return InputData{}, errors.New("end date is required")
	}
	endDateTime, err := time.Parse("2006-01-02", i.EndDate)
	if err != nil {
		return InputData{}, errors.New("end date invalid format")
	}

	if startDateTime.After(endDateTime) {
		return InputData{}, errors.New("start date must be before end date")
	}

	return InputData{
		SystemFile: i.SystemFile,
		BankFiles:  bankFilesList,
		StartDate:  startDateTime,
		EndDate:    endDateTime,
	}, nil
}
