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
	"strings"
	"text/tabwriter"
)

var (
	apiKey     string
	startDate  string
	endDate    string
	employeeId int
)

const (
	ActionList = "list"
	ActionAdd  = "add"
)

var actions = []string{ActionAdd, ActionList}

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

type Report struct {
	days           map[string]DayReport
	totalWorkHours float64
}
type DayReport struct {
	workHours float64
}

func init() {
	flag.StringVar(&apiKey, "apiKey", "", "Your BambooHR API key")
	flag.IntVar(&employeeId, "employeeId", 0, "Your BambooHR employee ID")
	flag.StringVar(&startDate, "start", "", "Start date filter for tracked working hours")
	flag.StringVar(&endDate, "end", "", "End date filter for tracked working hours")

	flag.Parse()
}

func main() {
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

	action := flag.Arg(0)
	switch action {
	case ActionList:
		processList()
		os.Exit(0)
	case ActionAdd:
		fmt.Println("selected add")
		os.Exit(0)
	default:
		fmt.Printf("No argument provided. You need to choose one of the supported actions: %s \n", strings.Join(actions, ", "))
		os.Exit(1)
	}
}

func processList() {
	workingHours, err := fetchWorkingHours()
	if err != nil {
		fmt.Printf("Failed fetching working hours: %v \n", err)
		os.Exit(1)
	}
	report := groupHoursByDate(workingHours)

	w := tabwriter.NewWriter(os.Stdout, 0, 5, 5, ' ', 0)
	// table header
	fmt.Fprintf(w, "Date\tTotal\t\n")
	for date, report := range report.days {
		fmt.Fprintf(w, "%s\t%s\n", date, convertDecimalTimeToTime(report.workHours))
	}

	fmt.Fprintf(w, "\nYour total working hours: %s \n", convertDecimalTimeToTime(report.totalWorkHours))
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

func groupHoursByDate(workingHours []TimeEntry) Report {
	dateMap := make(map[string]DayReport)
	totalHours := 0.0

	for _, entry := range workingHours {
		dayReport, ok := dateMap[entry.Date]
		if !ok {
			dateMap[entry.Date] = DayReport{workHours: entry.Hours}
			totalHours += entry.Hours
			continue
		}
		dayReport.workHours += entry.Hours
		dateMap[entry.Date] = dayReport
		totalHours += entry.Hours
	}

	return Report{
		dateMap,
		totalHours,
	}
}

func convertDecimalTimeToTime(decimalTime float64) string {
	hours := int(decimalTime)
	minutes := int(math.Round((decimalTime - float64(hours)) * 60))

	return fmt.Sprintf("%d hours and %d minutes", hours, minutes)
}
