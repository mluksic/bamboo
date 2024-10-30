package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"time"
)

func loadHolidays() (map[string]string, error) {
	var filename = "slovenian_public_work_off_days.csv"
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to open file: %v \n", err))
	}

	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ';'

	// buf read row by row
	holidays := make(map[string]string)

	// skip header row
	_, err = r.Read()
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
		isDayOff := row[3]
		notDayOffStrings := []string{"ne", "ne "}
		if slices.Contains(notDayOffStrings, isDayOff) {
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
