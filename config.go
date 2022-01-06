package main

import (
	"encoding/json"
	"fmt"
	"garage_counter/logger"
	"io/ioutil"
)

type Database struct {
	Ip string `json:"ip"`
	Port string `json:"port"`
	Database string `json:"database"`
	User string `json:"user"`
	Password string `json:"password"`
}

var dbCon Database

func ReadConfig() error {
	fmt.Println("Reading config file...")

	file, err := ioutil.ReadFile("@extensions/grc_config.json")

	if err != nil {
		logger.PrintLog(err.Error())
		return err
	}

	err = json.Unmarshal(file, &dbCon)

	if err != nil {
		logger.PrintLog(err.Error())
		return err
	}

	return nil
}