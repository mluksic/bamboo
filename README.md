# Bamboo CLI

CLI tool for automating the generation of daily work hour entries

## Features
- Automatically generates work hour entries for each workday between a given start and end date.
- Configurable daily work hours with three time slots: morning, break, and afternoon.
- Randomized work hours to simulate realistic entries.
- Excludes weekends, public holidays, and specified PTO days.
- Populates missing work entries on BambooHR

## Prerequisites

Download and install:

- [Go](https://go.dev/doc/install) (1.23+ recommended)


## Configuration

### Config File
The application uses a configuration file, `config.json`, which stores default values for your Bamboo `apiToken` and `employeeId`
```json
{
    "apiToken": "yourBambooApiToken",
    "employeeId": 123
}
```

## Building the app

```bash
$ go build -o bamboo main.go
```

## Running the app

### `list` command (skip config params if they're stored in `config.json`)
```bash
$ ./bamboo --apiKey yourBambooApiToken --employeeId 123  --start 2024-09-01 --end 2024-10-01 list
```

### `add` command (skip config params if they're stored in `config.json`)
```bash
$ ./bamboo --apiKey yourBambooApiToken --employeeId 123  --start 2024-09-01 --end 2024-10-01 --excludeDays 2024-09-15,2024-09-20 add
```

## Options
- `--apiKey` (**Required**) API token for BambooHR authentication
- `--employeeId`: (**Required**) Employee ID for whom the entries are generated - found in your BambooHR's URL
- `--start`: (**Required**) Start date in YYYY-MM-DD format
- `--end`: (**Required**) End date in YYYY-MM-DD format
- `--excludeDays`: (__Optional__) Comma-separated list of PTO dates in YYYY-MM-DD format. These dates will be excluded from work hour entries

## Example
### Generate work entries for October 2024, excluding October 28th, 29th and October 30th for PTO, October 31st is public holiday

(skip config params if they're stored in `config.json`)
```bash
$ ./bamboo --apiKey yourBambooApiToken --employeeId 123 --start 2024-10-01 --end 2024-11-01 --excludeDays 2024-10-28,2024-10-29,2024-10-30
```

#### Response
```
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

### Show work hours for September 2024, using `list` command

(skip config params if they're stored in `config.json`)
```bash
$ ./bamboo --apiKey yourBambooApiToken --employeeId 123 --start 2024-09-01 --end 2024-10-01 list
```

#### Response
```
Date           Weekday       Total
2024-09-02     Monday        8 hours and 14 minutes
2024-09-03     Tuesday       8 hours and 10 minutes
2024-09-04     Wednesday     8 hours and 18 minutes
2024-09-05     Thursday      8 hours and 5 minutes
2024-09-06     Friday        8 hours and 8 minutes
2024-09-09     Monday        7 hours and 13 minutes
2024-09-10     Tuesday       7 hours and 33 minutes
2024-09-11     Wednesday     8 hours and 8 minutes
2024-09-12     Thursday      7 hours and 44 minutes
2024-09-16     Monday        8 hours and 10 minutes
2024-09-17     Tuesday       8 hours and 1 minutes
2024-09-19     Thursday      8 hours and 11 minutes
2024-09-20     Friday        8 hours and 0 minutes
2024-09-23     Monday        7 hours and 46 minutes
2024-09-24     Tuesday       7 hours and 36 minutes
2024-09-25     Wednesday     8 hours and 20 minutes
2024-09-26     Thursday      7 hours and 56 minutes
2024-09-27     Friday        4 hours and 50 minutes
2024-09-30     Monday        7 hours and 30 minutes

Your total working hours: 147 hours and 53 minutes
```

After running the `add` command, double-check work entries in your Bamboo account

## Authors
👤 [Miha Luksic](https://www.mihaluksic.com)
