package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
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

	totalHours := 0.0
	for _, entry := range workingHours {
		totalHours = entry.Hours + totalHours
	}
	fmt.Printf("Your total working hours: %.2f hours \n", totalHours)
}

func fetchWorkingHours() ([]TimeEntry, error) {
	var workingHours TimeEntries
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

	err = json.Unmarshal(body, &workingHours)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to unmarshal json to struct: %v", err))
	}

	return workingHours, nil
}

//bamboo add --day 2024-01-01 --token asdf
//- create time entry 9:10 - 11:25, 11:25 - 11:52, 11:52 - 17:12
//
//bamboo get --month 2024-10
//- get time entries for specific month by date + working hoursapiKey, employeeId, startDate, endDate
