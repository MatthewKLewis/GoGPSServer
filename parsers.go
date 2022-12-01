package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type PreciseGPSData struct {
	deviceimei          string
	latitude            float64
	longitude           float64
	altitude            int
	devicetime          string
	speed               float64
	Batterylevel        int
	casefile_id         string
	address             string
	positioningmode     string
	tz                  string
	offender_name       string
	offender_id         string
	loc_message_content string
}

func getIMEIFromAP00(msg string) string {
	//fmt.Println(msg)
	//fmt.Println(msg[6:21])
	return msg[6:21]
}

func getIMEIFromLK(msg string) string {
	//fmt.Println(msg)
	//fmt.Println(msg[4:19])
	return msg[4:19]
}

// IWAP01301122V0000.0000N00000.0000E000.0201934000.0003900009700003,310,260,46136,11869187,AP1|16:8d:db:65:af:0f|-41&AP2|5c:5b:35:01:cd:d1|-43&AP3|3a:22:e2:a3:e1:07|-45&AP4|70:7d:b9:e1:95:c4|-47&AP5|70:7d:b9:e1:95:ce|-51#
func getJSONFromAP01(msg string, deviceIMEI string) (PreciseGPSData, error) {
	fmt.Println(msg)
	retObj := PreciseGPSData{
		deviceimei:          deviceIMEI,
		latitude:            0.0,
		longitude:           0.0,
		altitude:            0,
		devicetime:          "",
		speed:               0.0,
		Batterylevel:        0,
		casefile_id:         "",
		address:             "",
		positioningmode:     "",
		tz:                  "",
		offender_name:       "",
		offender_id:         "",
		loc_message_content: "",
	}

	// if msg[12] != 'A' { // Valid packets have an 'A' (65) rather than a 'V' (86) at index 6
	// 	return retObj, errors.New("No LAT LON")
	// }

	latDeg, err := strconv.ParseFloat(msg[13:15], 64) //Range is smaller for latitudes
	latPoints, err := strconv.ParseFloat(msg[16:23], 64)
	lonDeg, err := strconv.ParseFloat(msg[23:26], 64) //Range is 1 digit larger for longitudes
	lonPoints, err := strconv.ParseFloat(msg[27:33], 64)
	if err != nil {
		fmt.Println(err.Error())
		return retObj, err
	}
	retObj.latitude = latDeg + latPoints/300
	retObj.longitude = lonDeg + lonPoints/300
	fmt.Println(latDeg, latPoints, lonDeg, lonPoints)

	rssis := strings.Split(msg, "|")[1:]
	fmt.Println(rssis)

	return retObj, nil
}

func getJSONFromAP10(msg string, deviceIMEI string) (PreciseGPSData, error) {
	fmt.Println(msg)
	retObj := PreciseGPSData{
		deviceimei:      deviceIMEI,
		latitude:        0.0,
		longitude:       0.0,
		altitude:        0,
		devicetime:      "new go parser",
		speed:           0.0,
		Batterylevel:    0,
		casefile_id:     "new go parser",
		address:         "new go parser",
		positioningmode: "new go parser",
		tz:              "new go parser",
		offender_name:   "new go parser",
		offender_id:     "new go parser",
	}

	if msg[6] != 'A' {
		return retObj, errors.New("No LAT LON")
	}
	return retObj, nil
}
