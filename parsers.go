package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type PreciseGPSData struct {
	Deviceimei        string  `json:"deviceimei"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Altitude          int     `json:"altitude"`
	DeviceTime        string  `json:"devicetime"`
	Speed             float64 `json:"speed"`
	BatteryLevel      int     `json:"Batterylevel"`
	CasefileId        string  `json:"casefile_id"`
	Address           string  `json:"address"`
	PositioningMode   string  `json:"positioningmode"`
	Tz                string  `json:"tz"`
	OffenderName      string  `json:"offender_name"`
	OffenderId        string  `json:"offender_id"`
	LocMessageContent string  `json:"loc_message_content"`
}

func getIMEIFromAP00(msg string) string {
	return msg[6:21]
}

func getIMEIFromLK(msg string) string {
	return msg[4:19]
}

// IWAP01301122V0000.0000N00000.0000E000.0201934000.0003900009700003,310,260,46136,11869187,AP1|16:8d:db:65:af:0f|-41&AP2|5c:5b:35:01:cd:d1|-43&AP3|3a:22:e2:a3:e1:07|-45&AP4|70:7d:b9:e1:95:c4|-47&AP5|70:7d:b9:e1:95:ce|-51#
func getJSONFromAP01(msg string, deviceIMEI string) (PreciseGPSData, error) {
	retObj := PreciseGPSData{
		Deviceimei:        deviceIMEI,
		Latitude:          0.0,
		Longitude:         0.0,
		Altitude:          0,
		DeviceTime:        "",
		Speed:             0.0,
		BatteryLevel:      0,
		CasefileId:        "",
		Address:           "",
		PositioningMode:   "",
		Tz:                "",
		OffenderName:      "",
		OffenderId:        "",
		LocMessageContent: "",
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
	retObj.Latitude = latDeg + latPoints/300
	retObj.Longitude = lonDeg + lonPoints/300

	rssis := strings.Split(msg, "|")[1:]
	for i := 0; i < len(rssis); i++ {
		fmt.Println(rssis[i])
		if i%2 == 0 {
			retObj.LocMessageContent += rssis[i] + ":"
		} else if i == len(rssis)-1 {
			retObj.LocMessageContent += strings.Split(rssis[i], "&")[0] + "#"
		} else {
			retObj.LocMessageContent += strings.Split(rssis[i], "&")[0] + ";"
		}
	}
	return retObj, nil
}

func getJSONFromAP10(msg string, deviceIMEI string) (PreciseGPSData, error) {
	retObj := PreciseGPSData{
		Deviceimei:        deviceIMEI,
		Latitude:          0.0,
		Longitude:         0.0,
		Altitude:          0,
		DeviceTime:        "",
		Speed:             0.0,
		BatteryLevel:      0,
		CasefileId:        "",
		Address:           "",
		PositioningMode:   "",
		Tz:                "",
		OffenderName:      "",
		OffenderId:        "",
		LocMessageContent: "",
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
	retObj.Latitude = latDeg + latPoints/300
	retObj.Longitude = lonDeg + lonPoints/300
	fmt.Println(latDeg, latPoints, lonDeg, lonPoints)

	rssis := strings.Split(msg, "|")[1:]
	for i := 0; i < len(rssis); i++ {
		fmt.Println(rssis[i])
		if i%2 == 0 {
			retObj.LocMessageContent += rssis[i] + ":"
		} else if i == len(rssis)-1 {
			retObj.LocMessageContent += strings.Split(rssis[i], "&")[0] + "#"
		} else {
			retObj.LocMessageContent += strings.Split(rssis[i], "&")[0] + ";"
		}
	}

	return retObj, nil
}

func getJSONFromCUSTOMER(msg string, deviceIMEI string) (PreciseGPSData, error) {
	retObj := PreciseGPSData{
		Deviceimei:        deviceIMEI,
		Latitude:          0.0,
		Longitude:         0.0,
		Altitude:          0,
		DeviceTime:        "",
		Speed:             0.0,
		BatteryLevel:      0,
		CasefileId:        "",
		Address:           "",
		PositioningMode:   "",
		Tz:                "",
		OffenderName:      "",
		OffenderId:        "",
		LocMessageContent: "",
	}

	if msg[6] != 'A' {
		return retObj, errors.New("No LAT LON")
	}
	return retObj, nil
}
