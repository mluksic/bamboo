package main

import (
	"encoding/csv"
	"reflect"
	"strings"
	"testing"
)

func TestReadHolidaysFile(t *testing.T) {
	records := `DATUM;IME_PRAZNIKA;DAN_V_TEDNU;DELA_PROST_DAN;DAN;MESEC;LETO
1.01.2024;novo leto;ponedeljek;da;1;1;2024
2.01.2024;novo leto test;torek;da;2;1;2024
3.01.2024;novo leto test 2;sreda;ne;3;1;2024
4.01.2024;novo leto test 3;cetrtek;ne ;4;1;2024`

	want := map[string]string{
		"2024-01-01": "novo leto",
		"2024-01-02": "novo leto test",
	}

	reader := strings.NewReader(records)
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	h := &CsvHolidayFetcher{
		filepath: "test",
	}
	got, err := h.readHolidaysFile(csvReader)
	if err != nil {
		t.Errorf("readHolidaysFile() = '%v' should not return error", err)
	}
	if len(got) != 2 {
		t.Errorf("readHolidaysFile() should return 2 holidays, got %d", len(got))
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("readHolidaysFile() = %q, want %q", got, want)
	}
}

func TestMalformedHolidaysFile(t *testing.T) {
	invalidRecords := `DATUM;IME_PRAZNIKA;DAN_V_TEDNU;DELA_PROST_DAN;DAN;MESEC;LETO
1.01.2024;novo leto;ponedeljek;da;1;1;2024
2.01.2024"sasdf;;";novo leto test;torek;da;2;1;2024
3.01.2024;novo leto test 2;sreda;ne;3;1;2024
4.01.2024;novo leto test 3;cetrtek;ne ,4,1,2024`

	reader := strings.NewReader(invalidRecords)
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	h := &CsvHolidayFetcher{
		filepath: "test",
	}
	got, err := h.readHolidaysFile(csvReader)
	if err == nil || got != nil {
		t.Errorf("readHolidaysFile() = '%v' should return error", got)
	}
}