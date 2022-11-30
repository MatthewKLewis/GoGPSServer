package main

import (
	"errors"
	"fmt"
	"strconv"
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
	fmt.Println(msg)
	var data = strings.Split(msg, "'")[0]
	return data[6:21]
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

	// if msg[12] != 'A' { // Valid packets have an 'A' (65) rather than a 'V' (86) at index 6
	// 	return retObj, errors.New("No LAT LON")
	// }

	lat, _ := strconv.ParseFloat(msg[13:22], 64) //Range is smaller for latitudes
	lon, _ := strconv.ParseFloat(msg[23:33], 64) //Range is 1 digit larger for longitudes

	if msg[23] == 'S' {
		lat *= -1
	}
	if msg[23] == 'E' {
		lon *= -1
	}

	fmt.Println(lat)
	fmt.Println(lon)

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

func convertLatLon(latOrLon float64) {
	fmt.Println("Convert: ", latOrLon)
}
