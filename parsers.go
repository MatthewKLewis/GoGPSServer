package main

import (
	"fmt"
	"strconv"
	"strings"
)

// 39.52446611596893, -76.65204622381573
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
	Sos               bool    `json:"sos"`
}

func getIMEIFromAP00(msg string) string {
	return msg[6:21]
}

func getIMEIFromLK(msg string) string {
	fmt.Println("-LK-")
	fmt.Println(msg)
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
		Sos:               false,
	}

	if msg[12] != 'A' { // Valid packets have an 'A' (65) rather than a 'V' (86) at index 6
		retObj.Latitude = 39.52446611596893
		retObj.Longitude = -76.65204622381573
	} else {
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
		Sos:               false,
	}

	if msg[12] != 'A' { // Valid packets have an 'A' (65) rather than a 'V' (86) at index 6
		retObj.Latitude = 39.52446611596893
		retObj.Longitude = -76.65204622381573
	} else {
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

// [3G*357593065573357*0177*UDCUSTOMER1,011222,185310,V,0.0,N,0.0,E,22.0,0,-1,0,100,91,8382,0,00000000,1,1,310,260,46136,11869197,100,5,AiristaTesting,fa:55:3d:c0:32:4e,-47,,f6:55:3d:c0:32:4e,-22:13:78:5A:AF,-84,null,30:0A:B7:2D:7B:39,-52,null,02:37:01:E6:58:1E,-73,null,59:A3:33:97:9D:7A,-82,null,6D:D2:52:DC:1C:AB,-86,0.0]

// [3G*357593065573357*0107*ALCUSTOMER1,011222,200026,V,0.0,N,0.0,E,22.0,0,-1,0,100,86,8382,0,00100008,1,1,310,260,46136,11869197,100,5,AiristaTesting,fa:55:3d:c0:32:4e,-50,AiristaMist,5c:5b:35:01:b6:f1,-57,,ac:a3:1e:94:91:20,-68,Flex Point-2,70:03:7e:76:b4:3e,-71,,ba:3f:8c:fe:b9:85,-71,0.0]

// [3G*357593065573357*0107*ALCUSTOMER1,011222,200041,V,0.0,N,0.0,E,22.0,0,-1,0,100,86,8382,0,40000008,1,1,310,260,46136,11869197,100,5,AiristaTesting,fa:55:3d:c0:32:4e,-50,AiristaMist,5c:5b:35:01:b6:f1,-57,,ac:a3:1e:94:91:20,-68,Flex Point-2,70:03:7e:76:b4:3e,-71,,ba:3f:8c:fe:b9:85,-71,0.0]
// Missed call?

// [3G*357593065573357*0106*ALCUSTOMER1,011222,200140,V,0.0,N,0.0,E,22.0,0,-1,0,100,86,8382,0,00010000,1,1,310,260,46136,11869197,100,5,AiristaTesting,fa:55:3d:c0:32:4e,-45,Aruba,ac:a3:1e:94:91:21,-59,,ac:a3:1e:94:91:20,-59,AFDemo,a0:3d:6f:53:d6:84,-67,AiristaMist,5c:5b:35:01:b6:f1,-73,0.0]
// SOS?
func getJSONFromCUSTOMER(msg string, deviceIMEI string) (PreciseGPSData, error) {
	fmt.Println("-CUSTOMER-")
	fmt.Println(msg)
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
		Sos:               false,
	}
	splitString := strings.Split(msg, ",")

	lat, err := strconv.ParseFloat(splitString[4], 64)
	lon, err := strconv.ParseFloat(splitString[6], 64)
	handleError(err)

	retObj.Latitude = lat
	retObj.Longitude = lon

	// TEST DATA
	retObj.Latitude = 39.52446611596893
	retObj.Longitude = -76.65204622381573

	fmt.Println(splitString[16])

	splitStatuses := strings.Split(splitString[16], "")
	fmt.Println(splitStatuses)

	if splitStatuses[3] == "1" {
		fmt.Println("SOS BUTTON PRESSED!")
		retObj.Sos = true
	}

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

	// number of Base Stations := splitString[17]

	fmt.Println(retObj)
	return retObj, nil
}
