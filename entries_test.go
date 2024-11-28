package main

import (
	"reflect"
	"slices"
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
	tests := []struct {
		name string
		year int
		want YearReport
	}{
		{
			"Year2024Feb",
			2024,
			YearReport{
				month: map[string]MonthReport{"2024-02": {
					workDays:          20,
					holidays:          1,
					workHours:         160,
					totalHolidayHours: 8,
					totalHours:        168,
				}},
			},
		},
		{
			"Year2024Dec",
			2024,
			YearReport{
				month: map[string]MonthReport{"2024-12": {
					workDays:          20,
					holidays:          2,
					workHours:         160,
					totalHolidayHours: 16,
					totalHours:        176,
				}},
			},
		},
		{
			"Year2025Feb",
			2025,
			YearReport{
				month: map[string]MonthReport{"2025-02": {
					workDays:          20,
					holidays:          0,
					workHours:         160,
					totalHolidayHours: 0,
					totalHours:        160,
				}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := getRequiredHours(test.year)

			if len(got.month) < len(test.want.month) {
				t.Errorf("getRequiredHours() should generate entries for all months in year %d, got = %d", test.year, len(got.month))
			}

			months := getMonthsInMap(got.month)

			month, ok := got.month[months[0]]
			if !ok {
				t.Errorf("getRequiredHours() should have month %s in map %s", months[0], months[0])
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

	return months
}
