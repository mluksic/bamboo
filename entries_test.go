package main

import (
	"slices"
	"testing"
)

func TestGenerateWorkEntries(t *testing.T) {
	type args struct {
		report    Report
		startDate string
		endDate   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"ValidWeek",
			args{
				Report{
					map[string]DayReport{"2024-11-05": {7.7}},
					7.7,
				},
				"2024-10-25",
				"2024-11-06",
			},
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := generateWorkEntries(test.args.report, test.args.startDate, test.args.endDate)

			if (err != nil) != test.wantErr {
				t.Errorf("generateWorkEntries() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if len(got) != 7*3 {
				t.Errorf("generateWorkEntries() should have exactly 21 work entries (3x for each day, 7 work days in total), got %d", len(got))
				return
			}
			weekendDays := []string{"2024-10-26", "2024-10-27", "2024-11-02", "2024-11-03"}
			for _, entry := range got {
				if slices.Contains(weekendDays, entry.Date) {
					t.Errorf("generateWorkEntries() should not include %s (weekend)", entry.Date)
				}
				if entry.Date == test.args.endDate {
					t.Errorf("generateWorkEntries() should not include %s (end date)", test.args.endDate)
				}
				if _, ok := test.args.report.days[entry.Date]; ok {
					t.Errorf("generateWorkEntries() should not include entries for already populated date %s", entry.Date)
				}
			}
		})
	}
}
