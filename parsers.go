package main

import (
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
	BatteryLevel      int64   `json:"batterylevel"`
	CasefileId        string  `json:"casefile_id"`
	Address           string  `json:"address"`
	PositioningMode   string  `json:"positioningmode"`
	Tz                string  `json:"tz"`
	OffenderName      string  `json:"offender_name"`
	OffenderId        string  `json:"offender_id"`
	LocMessageContent string  `json:"loc_message_content"`
	Sos               bool    `json:"sos"`
}

type XpertDeviceData struct {
	Id int `json:"Id"`
}

type XpertConfigData struct {
	Id int `json:"Id"`
}

func getIMEIFromAP00(msg string) string {
	//fmt.Println(string(msg[6:21]))
	return msg[6:21]
}

func getIMEIFromLK(msg string) string {
	//fmt.Println(string(msg[4:19]))
	return msg[4:19]
}

// AP01221205V0000.0000N00000.0000E000.0175118000.0002800008000003,310,260,46136,228794903,AP1|fa:55:3d:c0:32:4e|-30&AP2|a0:3d:6f:53:d6:84|-38&AP3|a0:3d:6f:53:d6:8e|-39&AP4|a0:3d:6f:60:e8:70|-41&AP5|a0:3d:6f:60:e8:74|-42#
// IWAP01080524A2232.9806N11404.9355E000.1061830323.8706000908000102,460,0,9520,3671,Home|74-DE-2B-44-88-8C|97& Home1|74-DE-2B-44-88-8C|97&Home2|74-DE-2B-44-88-8C|97& Home3|74-DE-2B-44-88-8C|97#
// IWAP01221205V0000.0000N00000.0000E000.0191300000.0003200006000003,310,260,46136,11869197,AP1|a0:3d:6f:53:d6:8e|-30&AP2|a0:3d:6f:53:d6:84|-31&AP3|a0:3d:6f:53:d6:80|-32&AP4|5c:5b:35:01:b6:f1|-34&AP5|fa:55:3d:c0:32:4e|-35#
// 39.14554281096903, -76.87661699902404
func getJSONFromAP01(msg string, deviceIMEI string) (PreciseGPSData, error) {
	retObj := PreciseGPSData{
		Deviceimei:        deviceIMEI,
		Latitude:          0.0,
		Longitude:         0.0,
		Altitude:          0,
		DeviceTime:        "",
		Speed:             0.0,
		BatteryLevel:      100,
		CasefileId:        "",
		Address:           "",
		PositioningMode:   "",
		Tz:                "",
		OffenderName:      "",
		OffenderId:        "",
		LocMessageContent: "",
		Sos:               false,
	}

	if msg[12] != 'A' { // Valid packets have an 'A' (65) rather than a 'V' (86) at index 6
		retObj.Latitude = 0
		retObj.Longitude = 0
		// err := errors.New("No Lat Lon")
		// return retObj, err
	} else {
		latDeg, err := strconv.ParseFloat(msg[13:15], 64) //Range is smaller for latitudes
		latMins, err := strconv.ParseFloat(msg[15:18], 64)
		latSecs, err := strconv.ParseFloat(msg[18:22], 64)
		latSign := msg[22]

		lonDeg, err := strconv.ParseFloat(msg[23:26], 64) //Range is 1 digit larger for longitudes
		lonMins, err := strconv.ParseFloat(msg[26:29], 64)
		lonSecs, err := strconv.ParseFloat(msg[29:33], 64)
		lonSign := msg[33]

		fmt.Println(latDeg, latMins, latSecs, string(latSign), lonDeg, lonMins, lonSecs, string(lonSign))
		if err != nil {
			fmt.Println(err.Error())
			return retObj, err
		}
		retObj.Latitude = latDeg + latMins/60 + latSecs/360000
		retObj.Longitude = lonDeg + lonMins/60 + lonSecs/360000

		if string(latSign) == "S" {
			retObj.Latitude *= -1
		}
		if string(lonSign) == "W" {
			retObj.Longitude *= -1
		}
	}

	headers := strings.Split(msg, ",")[0]
	fmt.Println("Battery:", headers[57:60])
	if string(headers[57]) == "0" {
		retObj.BatteryLevel, _ = strconv.ParseInt(string(headers[58:60]), 0, 32)
	} else {
		retObj.BatteryLevel, _ = strconv.ParseInt(string(headers[57:60]), 0, 32)
	}

	rssis := strings.Split(msg, "|")[1:]
	for i := 0; i < len(rssis); i++ {
		if i%2 == 0 {
			retObj.LocMessageContent += rssis[i] + ":"
		} else if i == len(rssis)-1 {
			retObj.LocMessageContent += strings.Split(rssis[i], "&")[0]
		} else {
			retObj.LocMessageContent += strings.Split(rssis[i], "&")[0] + ";"
		}
	}
	return retObj, nil
}

func getJSONFromAP10(msg string, deviceIMEI string) (PreciseGPSData, error) {
	//fmt.Println(msg)
	retObj := PreciseGPSData{
		Deviceimei:        deviceIMEI,
		Latitude:          0.0,
		Longitude:         0.0,
		Altitude:          0,
		DeviceTime:        "",
		Speed:             0.0,
		BatteryLevel:      100,
		CasefileId:        "",
		Address:           "",
		PositioningMode:   "",
		Tz:                "",
		OffenderName:      "",
		OffenderId:        "",
		LocMessageContent: "",
		Sos:               false,
	}

	if msg[12] != 'A' { // Valid packets have an 'A' (65) rather than a 'V' (86) at index 6
		retObj.Latitude = 0
		retObj.Longitude = 0
		// err := errors.New("No Lat Lon")
		// return retObj, err
	} else {
		latDeg, err := strconv.ParseFloat(msg[13:15], 64) //Range is smaller for latitudes
		latMins, err := strconv.ParseFloat(msg[15:18], 64)
		latSecs, err := strconv.ParseFloat(msg[18:22], 64)
		latSign := msg[22]

		lonDeg, err := strconv.ParseFloat(msg[23:26], 64) //Range is 1 digit larger for longitudes
		lonMins, err := strconv.ParseFloat(msg[26:29], 64)
		lonSecs, err := strconv.ParseFloat(msg[29:33], 64)
		lonSign := msg[33]

		fmt.Println(latDeg, latMins, latSecs, string(latSign), lonDeg, lonMins, lonSecs, string(lonSign))
		if err != nil {
			fmt.Println(err.Error())
			return retObj, err
		}
		retObj.Latitude = latDeg + latMins/60 + latSecs/3600
		retObj.Longitude = lonDeg + lonMins/60 + lonSecs/3600

		if string(latSign) == "S" {
			retObj.Latitude *= -1
		}
		if string(lonSign) == "W" {
			retObj.Longitude *= -1
		}
	}

	// SOS and Band Information
	headers := strings.Split(msg, ",")
	if string(headers[5]) == "01" {
		retObj.Sos = true
	}

	firstBlock := headers[0]
	retObj.BatteryLevel, _ = strconv.ParseInt(firstBlock[57:60], 0, 32)

	// RSSI Information
	rssis := strings.Split(msg, "|")[1:]
	for i := 0; i < len(rssis); i++ {
		if i%2 == 0 {
			retObj.LocMessageContent += rssis[i] + ":"
		} else if i == len(rssis)-1 {
			retObj.LocMessageContent += strings.Split(rssis[i], "&")[0]
		} else {
			retObj.LocMessageContent += strings.Split(rssis[i], "&")[0] + ";"
		}
	}

	return retObj, nil
}

func getJSONFromCUSTOMER(msg string, deviceIMEI string) (PreciseGPSData, error) {
	//fmt.Println(msg)
	retObj := PreciseGPSData{
		Deviceimei:        deviceIMEI,
		Latitude:          0.0,
		Longitude:         0.0,
		Altitude:          0,
		DeviceTime:        "",
		Speed:             0.0,
		BatteryLevel:      100,
		CasefileId:        "",
		Address:           "",
		PositioningMode:   "",
		Tz:                "",
		OffenderName:      "",
		OffenderId:        "",
		LocMessageContent: "",
		Sos:               false,
	}
	splitString := strings.Split(msg, ",")

	lat, err := strconv.ParseFloat(splitString[4], 64)
	lon, err := strconv.ParseFloat(splitString[6], 64)
	handleError(err)

	retObj.Latitude = lat
	retObj.Longitude = lon

	retObj.BatteryLevel, _ = strconv.ParseInt(splitString[13], 0, 32)

	// TEST DATA
	retObj.Latitude = 39.52446611596893
	retObj.Longitude = -76.65204622381573

	if splitString[16] == "00010000" {
		retObj.Sos = true
	}

	//BASE STATIONS
	//numberofBaseStations, _ := strconv.ParseInt(splitString[17], 0, 32) //for each number here, there will be 6 base station info commas 17:24, 24:33, 33:40 etc...

	//WIFI SIGNALS
	indexOfNumberOfWifiSignals := 24 //or 17 + (number of base stations * 7?)
	numberOfWifiSignals, _ := strconv.ParseInt(splitString[indexOfNumberOfWifiSignals], 0, 32)
	indexOfMAC := indexOfNumberOfWifiSignals + 2

	//fmt.Println(indexOfNumberOfWifiSignals, numberOfWifiSignals)

	for i := 0; i < int(numberOfWifiSignals); i++ {
		retObj.LocMessageContent += splitString[int(indexOfMAC)] + ":"
		retObj.LocMessageContent += splitString[int(indexOfMAC)+1]
		indexOfMAC += 3

		if i == int(numberOfWifiSignals)-1 {
			retObj.LocMessageContent += "#"
		} else {
			retObj.LocMessageContent += ";"
		}
	}
	return retObj, nil
}

//statuses := splitString[16]
// fourthHex, _ := strconv.ParseInt(statuses[0:2], 0, 32)
// thirdHex, _ := strconv.ParseInt(statuses[2:4], 0, 32)
// secondHex, _ := strconv.ParseInt(statuses[4:6], 0, 32)
// firstHex, _ := strconv.ParseInt(statuses[6:8], 0, 32)
// fourthBits := fmt.Sprintf("%b", fourthHex)
// thirdBits := fmt.Sprintf("%b", thirdHex)
// secondBits := fmt.Sprintf("%b", secondHex)
// firstBits := fmt.Sprintf("%b", firstHex)
// fmt.Println(fourthBits, thirdBits, secondBits, firstBits)

// header := splitString[0]
// date := splitString[1]
// time := splitString[2]
// valid := splitString[3]
// latMark := splitString[5]
// lonMark := splitString[7]
// velocity := splitString[6]
// direction := splitString[7]
// altitude := splitString[8]
// numberOfSatellites := splitString[9]
// gsmIntensity := splitString[10]
// batteryPercent := splitString[13]
// stepCount := splitString[12]
// tumblingTimes := splitString[13]
// statusBinaries := splitString[16]
