package main

import (
	"embed"
	"encoding/csv"
	"encoding/json"
	"errors"
	fmt "fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

//go:embed slovenian_public_work_off_days.csv
var holidayFile embed.FS

type HolidayFetcher interface {
	loadHolidays() (map[string]string, error)
	fetchTimeOff() (map[string]string, error)
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

	var timeOffs = make(map[string]string)
	// skip fetching time offs if employeeID is not set
	if employeeId > 0 {
		timeOffs, err = h.fetchTimeOff()
		if err != nil {
			return nil, err
		}
	}
	holidays, err := h.readHolidaysFile(r)
	if err != nil {
		return nil, err
	}

	// combine public holidays and time offs
	outs := make(map[string]string)
	if timeOffs != nil {
		for k, v := range timeOffs {
			outs[k] = v
		}
	}

	for k, v := range holidays {
		outs[k] = v
	}

	return outs, nil
}

func (h *CsvHolidayFetcher) fetchTimeOff() (map[string]string, error) {
	var urlTemplate = "https://%s:x@api.bamboohr.com/api/gateway.php/flaviar/v1/time_off/whos_out?start=%s&end=%s"
	url := fmt.Sprintf(urlTemplate, apiKey, startDate, endDate)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to create whos out request: %v \n", err))
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to fetch whos out: %v \n", err))
	}
	defer resp.Body.Close()

	var resJson []struct {
		Id         int    `json:"id"`
		Type       string `json:"type"`
		EmployeeId int    `json:"employeeId"`
		Name       string `json:"name"`
		Start      string `json:"start"`
		End        string `json:"end"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read body: %v \n", err))
	}
	if err := json.Unmarshal(body, &resJson); err != nil {
		return nil, errors.New(fmt.Sprintf("unable to marshal response: %v \n", err))
	}

	outDays := make(map[string]string)
	for _, outEntry := range resJson {
		if outEntry.EmployeeId != employeeId {
			continue
		}

		entryStart, err := time.Parse("2006-01-02", outEntry.Start)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("unable to parse start time: %v \n", err))
		}
		entryEnd, err := time.Parse("2006-01-02", outEntry.End)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("unable to parse end time: %v \n", err))
		}

		if entryEnd == entryStart {
			outDays[outEntry.Start] = "timeOff"
		}
		if entryStart.After(entryEnd) {
			return nil, errors.New(fmt.Sprintf("entry start %s should not be after entry end %s \n", outEntry.Start, outEntry.End))
		}

		for d := entryStart; !d.After(entryEnd); d = d.AddDate(0, 0, 1) {
			outDays[d.Format("2006-01-02")] = "timeOff"
		}
	}

	return outDays, nil
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
