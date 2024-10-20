package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
)

var (
	apiKey     string
	startDate  string
	endDate    string
	employeeId int
)

type TimeEntries []TimeEntry

type TimeEntry struct {
	Id          int
	EmployeeId  int
	Type        string
	Date        string
	Start       string
	End         string
	Timezone    string
	Hours       float64
	Note        string
	ProjectInfo string
	ApprovedAt  string
	Approved    bool
}

type DayReport struct {
	workHours float64
}

func main() {
	flag.StringVar(&apiKey, "apiKey", "", "Your BambooHR api key for API access")
	flag.IntVar(&employeeId, "employeeId", 0, "Your BambooHr employee ID")
	flag.StringVar(&startDate, "start", "", "Start date filter for tracked working hours")
	flag.StringVar(&endDate, "end", "", "End date filter for tracked working hours")

	flag.Parse()

	if apiKey == "" {
		fmt.Println("Invalid 'apiKey' provided. Aborting")
		os.Exit(1)
	}
	if employeeId == 0 {
		fmt.Println("Invalid 'employeeId' provided. Aborting")
		os.Exit(1)
	}
	if startDate == "" {
		fmt.Println("Invalid 'start' date filter provided. Aborting")
		os.Exit(1)
	}
	if endDate == "" {
		fmt.Println("Invalid 'end' date filter provided. Aborting")
		os.Exit(1)
	}

	workingHours, err := fetchWorkingHours()
	if err != nil {
		fmt.Printf("Failed fetching working hours: %v \n", err)
		os.Exit(1)
	}

	dateMap := make(map[string]DayReport)
	for _, entry := range workingHours {
		dayReport, ok := dateMap[entry.Date]
		fmt.Println(entry.Hours)
		if !ok {
			dateMap[entry.Date] = DayReport{workHours: entry.Hours}
			continue
		} else {
			dayReport.workHours += entry.Hours
			dateMap[entry.Date] = dayReport
		}
	}

	for date, report := range dateMap {
		fmt.Printf("%s: %s\n", date, convertDecimalTimeToTime(report.workHours))
	}

	totalHours := calculateWorkingHours(dateMap)
	fmt.Printf("Your total working hours: %s \n", convertDecimalTimeToTime(totalHours))
}

func fetchWorkingHours() ([]TimeEntry, error) {
	var getHoursUrlTemplate = "https://%s:x@api.bamboohr.com/api/gateway.php/flaviar/v1/time_tracking/timesheet_entries?employeeIds=%d&start=%s&end=%s"
	url := fmt.Sprintf(getHoursUrlTemplate, apiKey, employeeId, startDate, endDate)

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to get tracked working hours from Bamboo: %v \n", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read response body: %v", err))
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("invalid 'apiToken' provided - API returned 401 (Unauthorized). Aborting")
	}
	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.New(string(body))
	}

	var workingHours TimeEntries
	err = json.Unmarshal(body, &workingHours)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to unmarshal json to struct: %v", err))
	}

	return workingHours, nil
}

func calculateWorkingHours(dates map[string]DayReport) float64 {
	totalHours := 0.0

	for _, date := range dates {
		totalHours = date.workHours + totalHours
	}

	return totalHours
}

func convertDecimalTimeToTime(decimalTime float64) string {
	hours := int(decimalTime)
	minutes := int(math.Round((decimalTime - float64(hours)) * 60))

	return fmt.Sprintf("%d hours and %d minutes", hours, minutes)
}
