package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
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
	ActionList = "list"
	ActionAdd  = "add"
)

var actions = []string{ActionAdd, ActionList}

type TimeEntriesPostBody struct {
	Entries []Entry `json:"entries"`
}
type Entry struct {
	EmployeeId int    `json:"employeeId"`
	Date       string `json:"date"`
	Start      string `json:"start"`
	End        string `json:"end"`
}

type TimeEntries []TimeEntry

type TimeEntry struct {
	Id          int         `json:"id"`
	EmployeeId  int         `json:"employeeId"`
	Type        string      `json:"type"`
	Date        string      `json:"date"`
	Start       time.Time   `json:"start"`
	End         time.Time   `json:"end"`
	Timezone    string      `json:"timezone"`
	Hours       float64     `json:"hours"`
	Note        interface{} `json:"note"`
	ProjectInfo interface{} `json:"projectInfo"`
	ApprovedAt  time.Time   `json:"approvedAt"`
	Approved    bool        `json:"approved"`
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
	flag.StringVar(&excludeDays, "excludeDays", "", "Comma-separated list of days (YYYY-MM-DD,YYYY-MM-DD) eg PTO, Collective Leave etc.")

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

	holidays, err = loadHolidays()
	if err != nil {
		fmt.Printf("Cannot load holidays: %v . Aborting \n", err)
		os.Exit(2)
	}
	excludedDays = loadExcludedDays()
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
	default:
		fmt.Printf("No argument provided. You need to choose one of the supported actions: %s \n", strings.Join(actions, ", "))
		os.Exit(1)
	}
}

func loadExcludedDays() map[string]bool {
	excludedDays := make(map[string]bool)
	if excludeDays == "" {
		return excludedDays
	}
	dates := strings.Split(excludeDays, ",")
	for _, date := range dates {
		excludedDays[date] = true
	}

	return excludedDays
}

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

func processList(report Report) {
	w := tabwriter.NewWriter(os.Stdout, 0, 5, 5, ' ', 0)
	// table header
	fmt.Fprintf(w, "Date\tWeekday\tTotal\t\n")

	// sort dates in asc order because map sorting order is random
	dates := make([]string, 0, len(report.days))
	for date := range report.days {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	for _, date := range dates {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			fmt.Printf("Unable to parse date from string: %v \n", err)
			os.Exit(1)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", date, t.Weekday(), convertDecimalTimeToTime(report.days[date].workHours))
	}

	fmt.Fprintf(w, "\nYour total working hours: %s \n", convertDecimalTimeToTime(report.totalWorkHours))
}

func addWorkingHours(report Report) {
	var storeHoursUrlTemplate = "https://%s:x@api.bamboohr.com/api/gateway.php/flaviar/v1/time_tracking/clock_entries/store"
	url := fmt.Sprintf(storeHoursUrlTemplate, apiKey)

	entries, err := generateWorkEntries(report)
	if err != nil {
		fmt.Printf("Unable to create post request entries: %v", err)
		os.Exit(1)
	}

	body, _ := json.Marshal(TimeEntriesPostBody{Entries: entries})
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Unable to create POST request: %v \n", err)
		os.Exit(1)
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Unable to trigger POST request: %v \n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Unable to read response body: %v", err)
		os.Exit(1)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Printf("invalid 'apiToken' provided - API returned 401 (Unauthorized). Aborting \n")
		os.Exit(1)
	}
	if resp.StatusCode == http.StatusBadRequest {
		fmt.Printf("Received Bad request (400): %s \n", string(respBody))
		os.Exit(1)
	}

	fmt.Println("Successfully populated working hour entries between two dates")
}

func generateWorkEntries(report Report) ([]Entry, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to parse start date: %v \n", err))
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to parse start date: %v \n", err))
	}
	// stop if diff between start and end date is more than 31 days
	if end.Sub(start).Hours() > 31*24 {
		return nil, errors.New(fmt.Sprint("max diff between days is 31 days \n", err))
	}

	weekend := []time.Weekday{time.Saturday, time.Sunday}
	existingHours := report.days
	var entries []Entry
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for s := start; !s.After(end); s = s.AddDate(0, 0, 1) {
		// exclude end date
		if s == end {
			break
		}
		// exclude weekends
		if slices.Contains(weekend, s.Weekday()) {
			fmt.Printf("Excluded %s because it's a weekend \n", s.Format("2006-01-2"))
			continue
		}
		// exclude days when hours were already logged
		_, ok := existingHours[s.Format("2006-01-02")]
		if ok {
			fmt.Printf("Excluded '%s' because hours were already logged for this day \n", s.Format("2006-01-2"))
			continue
		}
		// exclude holidays
		holiday, ok := holidays[s.Format("2006-01-02")]
		if ok {
			fmt.Printf("Excluded '%s' because it's public holiday - %s \n", s.Format("2006-01-2"), holiday)
			continue
		}
		// exclude provided dates (PTOs, collective leave, etc.)
		_, ok = excludedDays[s.Format("2006-01-02")]
		if ok {
			fmt.Printf("Excluded '%s' because you excluded it \n", s.Format("2006-01-2"))
			continue
		}

		// 440 to 460 minutes (7h20 to 7h40) - excluding 30min for break
		workMinutes := r.Intn(21) + 440
		// randomly select either 8 or 8 as the hour, and random minute within the hour (8AM - 9:59AM)
		morningStart := time.Date(s.Year(), s.Month(), s.Day(), r.Intn(2)+8, r.Intn(60), 0, 0, s.Location())
		// calculate the halfway of the work duration
		morningEnd := morningStart.Add(time.Duration(workMinutes/2) * time.Minute)
		breakStart := morningEnd
		// 30min lunch break
		breakEnd := breakStart.Add(30 * time.Minute)
		afternoonEnd := breakEnd.Add(time.Duration(workMinutes-workMinutes/2) * time.Minute)

		startEntry := Entry{
			EmployeeId: employeeId,
			Date:       s.Format("2006-01-02"),
			Start:      morningStart.Format("15:04"),
			End:        morningEnd.Format("15:04"),
		}

		lunchEntry := Entry{
			EmployeeId: employeeId,
			Date:       s.Format("2006-01-02"),
			Start:      breakStart.Format("15:04"),
			End:        breakEnd.Format("15:04"),
		}

		endEntry := Entry{
			EmployeeId: employeeId,
			Date:       s.Format("2006-01-02"),
			Start:      breakEnd.Format("15:04"),
			End:        afternoonEnd.Format("15:04"),
		}

		entries = append(entries, startEntry, lunchEntry, endEntry)
	}

	return entries, nil
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
