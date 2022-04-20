package main

import (
	"encoding/json"
)

func struct2JSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func parseVehicleArray(start int, args []string) []string {
	var vehicles []string
	for i := start; i < len(args); i++ {
		vehicles = append(vehicles, args[i])
	}
	return vehicles
}
