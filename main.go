package main

import (
	"flag"
	"fmt"
	"os"
)

var getHoursUrlTemplate string = "https://%s:x@api.bamboohr.com/api/gateway.php/flaviar/v1/time_tracking/timesheet_entries?employeeIds=%d&start=%s&end=%s"
var (
	apiKey     string
	startDate  string
	endDate    string
	employeeId int
)

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
	fmt.Printf(getHoursUrlTemplate, apiKey, employeeId, startDate, endDate)
}

//bamboo add --day 2024-01-01 --token asdf
//- create time entry 9:10 - 11:25, 11:25 - 11:52, 11:52 - 17:12
//
//bamboo get --month 2024-10
//- get time entries for specific month by date + working hoursapiKey, employeeId, startDate, endDate
