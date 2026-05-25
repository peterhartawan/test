package entity

import (
	"reflect"
	"testing"
	"time"
)

func TestInputRequest_ParseAndValidate(t *testing.T) {
	startDateTime, _ := time.Parse("2006-01-02", "2026-01-01")
	endDateTime, _ := time.Parse("2006-01-02", "2026-01-07")

	tests := []struct {
		name    string
		fields  InputRequest
		want    InputData
		wantErr bool
	}{
		{
			name: "invalid input system file empty",
			fields: InputRequest{
				SystemFile: "",
				BankFiles:  "bank1.csv,bank2.csv",
				StartDate:  "2026-01-01",
				EndDate:    "2026-01-07",
			},
			want:    InputData{},
			wantErr: true,
		},
		{
			name: "invalid system file not csv",
			fields: InputRequest{
				SystemFile: "valid.txt",
				BankFiles:  "bank1.csv",
				StartDate:  "2026-01-01",
				EndDate:    "2026-01-07",
			},
			want:    InputData{},
			wantErr: true,
		},
		{
			name: "invalid system file not found",
			fields: InputRequest{
				SystemFile: "./../../testdata/system/not_found.csv",
				BankFiles:  "bank1.csv",
				StartDate:  "2026-01-01",
				EndDate:    "2026-01-07",
			},
			want:    InputData{},
			wantErr: true,
		},
		{
			name: "invalid bank file empty",
			fields: InputRequest{
				SystemFile: "./../../testdata/system/valid.csv",
				BankFiles:  "",
				StartDate:  "2026-01-01",
				EndDate:    "2026-01-07",
			},
			want:    InputData{},
			wantErr: true,
		},
		{
			name: "invalid bank file not csv",
			fields: InputRequest{
				SystemFile: "./../../testdata/system/valid.csv",
				BankFiles:  "bank1.txt",
				StartDate:  "2026-01-01",
				EndDate:    "2026-01-07",
			},
			want:    InputData{},
			wantErr: true,
		},
		{
			name: "invalid bank file not found",
			fields: InputRequest{
				SystemFile: "./../../testdata/system/valid.csv",
				BankFiles:  "bank1.csv",
				StartDate:  "2026-01-01",
				EndDate:    "2026-01-07",
			},
			want:    InputData{},
			wantErr: true,
		},
		{
			name: "invalid start date empty",
			fields: InputRequest{
				SystemFile: "./../../testdata/system/valid.csv",
				BankFiles:  "./../../testdata/bank/valid.csv",
				StartDate:  "",
				EndDate:    "2026-01-07",
			},
			want:    InputData{},
			wantErr: true,
		},
		{
			name: "invalid start date format",
			fields: InputRequest{
				SystemFile: "./../../testdata/system/valid.csv",
				BankFiles:  "./../../testdata/bank/valid.csv",
				StartDate:  "2026/01/01",
				EndDate:    "2026-01-07",
			},
			want:    InputData{},
			wantErr: true,
		},
		{
			name: "invalid end date empty",
			fields: InputRequest{
				SystemFile: "./../../testdata/system/valid.csv",
				BankFiles:  "./../../testdata/bank/valid.csv",
				StartDate:  "2026-01-01",
				EndDate:    "",
			},
			want:    InputData{},
			wantErr: true,
		},
		{
			name: "invalid end date format",
			fields: InputRequest{
				SystemFile: "./../../testdata/system/valid.csv",
				BankFiles:  "./../../testdata/bank/valid.csv",
				StartDate:  "2026-01-01",
				EndDate:    "2026/01/07",
			},
			want:    InputData{},
			wantErr: true,
		},
		{
			name: "invalid start date after end date",
			fields: InputRequest{
				SystemFile: "./../../testdata/system/valid.csv",
				BankFiles:  "./../../testdata/bank/valid.csv",
				StartDate:  "2026-01-10",
				EndDate:    "2026-01-07",
			},
			want:    InputData{},
			wantErr: true,
		},
		{
			name: "valid input",
			fields: InputRequest{
				SystemFile: "./../../testdata/system/valid.csv",
				BankFiles:  "./../../testdata/bank/valid.csv,./../../testdata/bank/malformed_row.csv",
				StartDate:  "2026-01-01",
				EndDate:    "2026-01-07",
			},
			want: InputData{
				SystemFile: "./../../testdata/system/valid.csv",
				BankFiles:  []string{"./../../testdata/bank/valid.csv", "./../../testdata/bank/malformed_row.csv"},
				StartDate:  startDateTime,
				EndDate:    endDateTime,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &tt.fields
			got, err := i.ParseAndValidate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("InputRequest.ParseAndValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InputRequest.ParseAndValidate() = %v, want %v", got, tt.want)
			}
		})
	}
}
