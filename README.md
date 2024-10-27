# Bamboo CLI

CLI tool for automating the generation of daily work hour entries

## Features
- Automatically generates work hour entries for each workday between a given start and end date.
- Configurable daily work hours with three time slots: morning, break, and afternoon.
- Randomized work hours to simulate realistic entries.
- Excludes weekends, public holidays, and specified PTO days.
- Populates missing work entries on BambooHR

## Dependencies

- [Go](https://go.dev/doc/install)

## Prerequisites

Download and install:

- [Go](https://go.dev/doc/install)

## Building the app

```bash
$ go build -o bamboo
```

## Options
- `--apiKey` (**Required**) API token for BambooHR authentication
- `--employeeId`: (**Required**) Employee ID for whom the entries are generated - found in your BambooHR's URL
- `--start`: (**Required**) Start date in YYYY-MM-DD format
- `--end`: (**Required**) End date in YYYY-MM-DD format
- `--excludeDays`: (__Optional__) Comma-separated list of PTO dates in YYYY-MM-DD format. These dates will be excluded from work hour entries

## Running the app
### `list` command
```bash
$ go run . --apiKey myApiToken --employeeId yourEmployeeId  --start 2024-09-01 --end 2024-10-01 list
```

### `add` command
```bash
$ go run . --apiKey myApiToken --employeeId yourEmployeeId  --start 2024-09-01 --end 2024-10-01 --excludeDays 2024-09-15,2024-09-20 add
```

## Example
Generate work entries for October 2024, excluding October 28th, 29th and October 30th for PTO, October 31st is public holiday
```bash
$ go run . --apiKey myApiToken --employeeId 12345 --start 2024-10-01 --end 2024-11-01 --excludeDays 2024-10-28,2024-10-29,2024-10-30
```

## Response
```bash
Excluded '2024-10-1' because hours were already logged for this day
Excluded '2024-10-2' because hours were already logged for this day
Excluded '2024-10-3' because hours were already logged for this day
Excluded '2024-10-4' because hours were already logged for this day
Excluded 2024-10-5 because it's a weekend
Excluded 2024-10-6 because it's a weekend
Excluded 2024-10-12 because it's a weekend
Excluded 2024-10-13 because it's a weekend
Excluded '2024-10-16' because hours were already logged for this day
Excluded '2024-10-17' because hours were already logged for this day
Excluded 2024-10-19 because it's a weekend
Excluded 2024-10-20 because it's a weekend
Excluded 2024-10-26 because it's a weekend
Excluded 2024-10-27 because it's a weekend
Excluded '2024-10-28' because you excluded it
Excluded '2024-10-29' because you excluded it
Excluded '2024-10-30' because you excluded it
Excluded '2024-10-31' because it's public holiday - dan reformacije
Successfully populated working hour entries between two dates
```

👤 [Miha Luksic](https://www.mihaluksic.com)
