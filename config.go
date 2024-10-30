package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	ApiToken   string `json:"apiToken"`
	EmployeeId int    `json:"employeeId"`
}

func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.ReadFile(filename)
	if err != nil {
		return config, errors.New(fmt.Sprintf("unable to read config file: %v \n", err))
	}

	return readConfigFile(file)
}

func readConfigFile(r []byte) (Config, error) {
	var config Config
	err := json.Unmarshal(r, &config)
	if err != nil {
		return config, errors.New(fmt.Sprintf("unable to unmarshal JSON from config file: %v \n", err))
	}

	return config, nil
}
