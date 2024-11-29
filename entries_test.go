package main

import (
	"reflect"
	"slices"
	"sort"
	"testing"
	"time"
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

func TestDaysInMonth(t *testing.T) {
	type args struct {
		month time.Month
		year  int
	}

	tests := []struct {
		name  string
		input args
		want  int
	}{
		{
			"DaysInFeb2024",
			args{
				time.February,
				2024,
			},
			29,
		},
		{
			"DaysInMarch2024",
			args{
				time.March,
				2024,
			},
			31,
		},
		{
			"DaysInMarch2024",
			args{
				time.March,
				2024,
			},
			31,
		},
		{
			"DaysInFeb2025",
			args{
				time.February,
				2025,
			},
			28,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := daysInMonth(test.input.month, test.input.year)

			if got != test.want {
				t.Errorf("daysInMonth() = got = %d ; want = %d", got, test.want)
			}
		})
	}
}

func TestGetRequiredHours(t *testing.T) {
	type args struct {
		year     int
		holidays map[string]string
	}

	tests := []struct {
		name  string
		input args
		want  YearReport
	}{
		{
			"Year2024Feb",
			args{2024, map[string]string{"2024-02-08": "Prešernov dan"}},
			YearReport{
				map[string]MonthReport{"2024-02": {
					20,
					1,
					160,
					8,
					168,
				}},
			},
		},
		{
			"Year2024Dec",
			args{2024, map[string]string{"2024-12-25": "božič", "2024-12-26": "dan samostojnosti"}},
			YearReport{
				map[string]MonthReport{"2024-12": {
					20,
					2,
					160,
					16,
					176,
				}},
			},
		},
		{
			"Year2025Feb",
			args{2025, map[string]string{"2025-02-08": "Prešernov dan"}},
			YearReport{
				map[string]MonthReport{"2025-02": {
					20,
					0,
					160,
					0,
					160,
				}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := getRequiredHours(test.input.year, test.input.holidays)

			if len(got.month) < len(test.want.month) {
				t.Errorf("getRequiredHours() should generate entries for all months in year %d, got = %d", test.input.year, len(got.month))
			}

			months := getMonthsInMap(test.want.month)
			month, ok := got.month[months[0]]
			if !ok {
				t.Errorf("getRequiredHours() should have month %s in map %v", months[0], got.month)
			}
			if !reflect.DeepEqual(test.want.month[months[0]], month) {
				t.Errorf("getRequiredHours() want = %v ; got = %v", test.want.month[months[0]], month)
			}

		})
	}
}

func getMonthsInMap(monthsMap map[string]MonthReport) []string {
	var months []string

	for month := range monthsMap {
		months = append(months, month)
	}

	sort.Strings(months)
	return months
}
