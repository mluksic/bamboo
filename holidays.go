package main

import (
	"embed"
	"encoding/csv"
	"errors"
	fmt "fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"time"
)

//go:embed slovenian_public_work_off_days.csv
var holidayFile embed.FS

type HolidayFetcher interface {
	loadHolidays() (map[string]string, error)
}

type CsvHolidayFetcher struct {
	filepath string
}

func NewCsvHolidays(filepath string) *CsvHolidayFetcher {
	return &CsvHolidayFetcher{
		filepath: filepath,
	}
}

func (h *CsvHolidayFetcher) loadHolidays() (map[string]string, error) {
	file, err := holidayFile.Open(h.filepath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to open file: %v \n", err))
	}

	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ';'

	return h.readHolidaysFile(r)
}

func (h *CsvHolidayFetcher) readHolidaysFile(r *csv.Reader) (map[string]string, error) {
	holidays := make(map[string]string)

	// skip header row
	_, err := r.Read()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to skip header row: %v \n", err))
	}

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New(fmt.Sprintf("unable to read row: %v \n", err))
		}

		// exclude 'working' public holidays
		if !h.isOffDay(row[3]) {
			continue
		}

		year, _ := strconv.Atoi(row[6])
		month, _ := strconv.Atoi(row[5])
		day, _ := strconv.Atoi(row[4])
		holidayName := row[1]
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

		holidays[date.Format("2006-01-02")] = holidayName

	}

	return holidays, nil
}

func (h *CsvHolidayFetcher) isOffDay(separator string) bool {
	dayOffStrings := []string{"da"}

	if slices.Contains(dayOffStrings, separator) {
		return true
	}

	return false
}

func loadExcludedDays(excludeDays string) (map[string]bool, error) {
	excludedDays := make(map[string]bool)

	if len(excludeDays) == 0 {
		return excludedDays, nil
	}

	lastChar := string(excludeDays[len(excludeDays)-1])
	if lastChar == "," {
		return excludedDays, errors.New(fmt.Sprint("'excludeDays' should not contain spaces \n"))
	}

	dates := strings.Split(excludeDays, ",")
	for _, date := range dates {
		excludedDays[date] = true
	}

	return excludedDays, nil
}
