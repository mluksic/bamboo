package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	apiKey       string
	startDate    string
	endDate      string
	excludeDays  string
	employeeId   int
	holidays     map[string]string
	excludedDays map[string]bool
)

const (
	ActionList     = "list"
	ActionAdd      = "add"
	ActionRequired = "required"
)

var actions = []string{ActionAdd, ActionList, ActionRequired}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Printf("Unable to load config file - %v. Aborting", err)
		os.Exit(1)
	}

	flag.StringVar(&apiKey, "apiKey", config.ApiToken, "Your BambooHR API key")
	flag.IntVar(&employeeId, "employeeId", config.EmployeeId, "Your BambooHR employee ID")
	flag.StringVar(&startDate, "start", "", "Start date filter for tracked working hours")
	flag.StringVar(&endDate, "end", "", "End date filter for tracked working hours")
	flag.StringVar(&excludeDays, "excludeDays", "", "Comma-separated list of days (YYYY-MM-DD,YYYY-MM-DD) eg PTO, Collective Leave etc.")

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
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		fmt.Println("Unable to parse 'start' date. Aborting")
		os.Exit(1)
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		fmt.Println("Unable to parse 'end' date. Aborting")
		os.Exit(1)
	}
	if start.After(end) {
		fmt.Println("'end' date cannot be before 'start' date")
		os.Exit(1)
	}

	holidayFetcher := NewCsvHolidays("slovenian_public_work_off_days.csv")
	holidays, err = holidayFetcher.loadHolidays()
	if err != nil {
		fmt.Printf("Cannot load holidays: %v . Aborting \n", err)
		os.Exit(1)
	}
	excludedDays, err = loadExcludedDays(excludeDays)
	if err != nil {
		fmt.Printf("Cannot parse excluded days: %v", err)
		os.Exit(1)
	}
	workingHours, err := fetchWorkingHours()
	if err != nil {
		fmt.Printf("Failed fetching working hours: %v \n", err)
		os.Exit(1)
	}

	action := flag.Arg(0)

	report := groupHoursByDate(workingHours)

	switch action {
	case ActionList:
		processList(report)
		os.Exit(0)
	case ActionAdd:
		addWorkingHours(report)
		os.Exit(0)
	case ActionRequired:
		processRequiredHours()
		os.Exit(0)
	default:
		fmt.Printf("No argument provided. You need to choose one of the supported actions: %s \n", strings.Join(actions, ", "))
		os.Exit(1)
	}
}
