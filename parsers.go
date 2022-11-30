package main

import (
	"fmt"
	"strings"
)

type PreciseGPSData struct {
	deviceimei      string
	latitude        float64
	longitude       float64
	altitude        int
	devicetime      string
	speed           float64
	Batterylevel    int
	casefile_id     string
	address         string
	positioningmode string
	tz              string
	offender_name   string
	offender_id     string
}

func getIMEIFromAP00(msg string) string {
	var data = strings.Split(msg, "'")[0]
	return data[6:21]
}

func getIMEIFromLK(msg string) string {
	var data = strings.Split(msg, "'")[0]
	return data[1:15]
}

func getJSONFromAP01(msg string) PreciseGPSData {
	retObj := PreciseGPSData{
		deviceimei:      "string",
		latitude:        0.0,
		longitude:       0.0,
		altitude:        0,
		devicetime:      "string",
		speed:           0.0,
		Batterylevel:    0,
		casefile_id:     "string",
		address:         "string",
		positioningmode: "string",
		tz:              "string",
		offender_name:   "string",
		offender_id:     "string",
	}
	return retObj
}

func getJSONFromAP10(msg string) PreciseGPSData {
	retObj := PreciseGPSData{
		deviceimei:      "string",
		latitude:        0.0,
		longitude:       0.0,
		altitude:        0,
		devicetime:      "string",
		speed:           0.0,
		Batterylevel:    0,
		casefile_id:     "string",
		address:         "string",
		positioningmode: "string",
		tz:              "string",
		offender_name:   "string",
		offender_id:     "string",
	}
	return retObj
}

func convertLatLon(latOrLon float64) {
	fmt.Println("Convert: ", latOrLon)
}
