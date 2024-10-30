package main

import (
	"reflect"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    Config
		wantErr bool
	}{
		{
			"ValidConfig",
			[]byte(`{"apiToken":"myApiToken","employeeId":1234}`),
			Config{
				ApiToken:   "myApiToken",
				EmployeeId: 1234,
			},
			false,
		},
		{
			"EmptyConfig",
			[]byte(`{"apiToken":"","employeeId":0}`),
			Config{
				ApiToken:   "",
				EmployeeId: 0,
			},
			false,
		},
		{
			"MalformedJson",
			[]byte(`{"apiTok""","employd:0}"`),
			Config{},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, err := readConfigFile(test.input)

			if (err != nil) != test.wantErr {
				t.Errorf(`readConfigFile(bytes) should return an error", got %q`, config)
				return
			}
			if !reflect.DeepEqual(config, test.want) {
				t.Errorf(`readConfigFile(bytes) should match %v, got %v`, test.want, config)
			}
		})
	}
}
