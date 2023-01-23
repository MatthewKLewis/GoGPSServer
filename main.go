package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	host                  = "172.31.51.82"
	port                  = "8000"
	connectionType        = "tcp"
	locationRoute         = "http://52.45.17.177:802/XpertRestApi/api/location_data" //http://52.45.17.177:802/XpertRestApi/api/location_data //http://3.212.201.170:802/XpertRestApi/api/location_data
	getDeviceByMacRoute   = "http://3.212.201.170:802/XpertRestApi/api/Device/GetByMacAddress?"
	getConfigPendingRoute = "http://3.212.201.170:802/XpertRestApi/api/Device/GetConfigPendingByDeviceId?"
	setConfigRoute        = "http://3.212.201.170:802/XpertRestApi/api/Device/SetDeviceConfigurations?"
)

func main() {
	socketServer, err := net.Listen(connectionType, host+":"+port)
	handleError(err)

	defer socketServer.Close()

	fmt.Println("Listening on", host, ":", port)
	for {
		connection, err := socketServer.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go threadedClientConnectionHandler(connection)
	}
}

func threadedClientConnectionHandler(connection net.Conn) {

	deviceIMEI := ""
	oDevice := XpertDeviceData{}
	packetData := PreciseGPSData{}

	for {
		buffer := make([]byte, 2048)
		dtg := time.Now().Format("01/30/2006 15:04:05")

		// Blocks here, waiting for next message, sets time after Reading
		mLen, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Breaking from While Loop, Couldn't Read Buffer!")
			break
		}

		message := string(buffer[:mLen])

		fmt.Println("Message: " + message)

		if deviceIMEI != "" {
			oDevice = getDeviceFromIMEI(deviceIMEI)

			if oDevice.PendingConfigId != 0 {
				go pushConfigurationToTag(connection, deviceIMEI, oDevice)
			}
		}

		// Sorting
		if strings.Contains(message, "AP00") { // AP00 = Connection
			deviceIMEI = getIMEIFromAP00(message)
			go connection.Write([]byte("IWBP00," + dtg + ",4#"))

		} else if strings.Contains(message, "AP03") { // AP03 = Heartbeat
			go connection.Write([]byte("IWBP03#"))

		} else if strings.Contains(message, "IWAPTQ") { // IWAPTQ = ???
			// do nothing
		} else if strings.Contains(message, "IWAP01") { // AP01 = Location
			packetData, err = getJSONFromAP01(message, deviceIMEI) //carve off the "IW"
			if err != nil {
				fmt.Println("Breaking from While Loop, No Lat Lon!")
				break
			}
			go sendLocationToAPI(packetData)

		} else if strings.Contains(message, "AP01") { // AP01 = Location
			packetData, err = getJSONFromAP01("IW"+message, deviceIMEI)
			if err != nil {
				fmt.Println("Breaking from While Loop, No Lat Lon!")
				break
			}
			go sendLocationToAPI(packetData)

		} else if strings.Contains(message, "AP10") { // AP10 = Alert
			packetData, err = getJSONFromAP10(message, deviceIMEI)
			handleError(err)
			go connection.Write([]byte("IWBP10#"))
			go sendLocationToAPI(packetData)

		} else if strings.Contains(message, "AP12") { // AP10 = Alert
			fmt.Println("RECEIVED NEW PHONE NUMBERS!")
			sendSatisfiedConfigToAPI(oDevice)
			oDevice.PendingConfigId = 0

		} else if strings.Contains(message, "AP33") { // AP10 = Alert
			fmt.Println("RECEIVED NEW WORKING MODE!")
			//tell API that the config change has been satisfied

		} else if strings.Contains(message, "LK") { // LK = Link
			deviceIMEI = getIMEIFromLK(message)
			fmt.Println("SENDING: [3G*" + deviceIMEI + "*0027*SOS,15712257714,15712257714,15712257714]")

			connection.Write([]byte("[3G*" + deviceIMEI + "*0002*LK]"))
			connection.Write([]byte("[3G*" + deviceIMEI + "*0009*UPLOAD,10]"))
			connection.Write([]byte("[3G*" + deviceIMEI + "*0027*SOS,15555555555,15555555555,15555555555]"))
		} else if strings.Contains(message, "ALCUSTOMER") { // ALCUSTOMER = Alarm
			packetData, err = getJSONFromCUSTOMER(message, deviceIMEI)
			handleError(err)
			sendLocationToAPI(packetData)
		} else if strings.Contains(message, "UDCUSTOMER") { // UDCUSTOMER = Location
			packetData, err = getJSONFromCUSTOMER(message, deviceIMEI)
			handleError(err)
			sendLocationToAPI(packetData)

			//RESPONSES -- THREAD ENDS
		} else if strings.Contains(message, "UPLOAD") { // UDCUSTOMER = Location
			fmt.Println("UPLOAD complete.")
		} else if strings.Contains(message, "SOS") { // UDCUSTOMER = Location
			fmt.Println("SOS # UPDATE complete.")
		} else {
			fmt.Println("Message? " + message)
		}
	}
}

func sendLocationToAPI(packet PreciseGPSData) {
	body, err := json.Marshal(packet)
	handleError(err)

	fmt.Println(string(body))

	res, err := http.Post(locationRoute, "application/json", bytes.NewBuffer(body))
	handleError(err)

	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("Location Update API response: %s\n", resBody)
}

// getConfigRoute = "http://3.212.201.170:802/XpertRestApi/api/Device/GetConfigByDeviceId?  DeviceId=1 &CustomerId=1"
// setConfigRoute = "http://3.212.201.170:802/XpertRestApi/api/Device/SetDeviceConfigurations?  CustomerId=2047 &ConfigID=1235 &PendingConfigID=0"
func pushConfigurationToTag(connection net.Conn, imei string, device XpertDeviceData) {
	client := &http.Client{}
	configData := XpertConfigData{}
	configDefData := XpertW9ConfigDefinitionData{}
	var basicAuth = "Basic " + "YWZhZG1pbjphZG1pbg==" //Basic YWZhZG1pbjphZG1pbg==

	//Request Config from Device Id
	formattedGetConfigRoute := getConfigPendingRoute + "DeviceId=" + fmt.Sprint(device.Id) + "&CustomerId=" + "2047"
	req, err := http.NewRequest("GET", formattedGetConfigRoute, nil)
	req.Header.Add("Authorization", basicAuth)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	//Decode and fill Config Data
	if err := json.NewDecoder(resp.Body).Decode(&configData); err != nil {
		log.Fatal("Error decoding json" + err.Error())
	}

	//If no pending Config, drop
	if configData.Id == 0 {
		fmt.Println("No Pending Command for", device.PendingConfigId)
		return
	}

	//Decode JSON-within-JSON
	json.Unmarshal([]byte(configData.ConfigDef), &configDefData)

	//Write Phone Command to Device
	sosPhoneCommand := "IWBP12," + imei + ",080835," + configDefData.PhoneNumber1 + "," + configDefData.PhoneNumber2 + "," + configDefData.PhoneNumber3 + "#"
	//fmt.Println(sosPhoneCommand)
	connection.Write([]byte(sosPhoneCommand))

	//Write Mode Command To Device //ADD SWITCH CASE FOR 1,2,3 NORMAL,POWERSAVE,EMERGENCY
	workingModeCommand := "IWBP33," + imei + ",080835," + "3" + "#"
	//fmt.Println(workingModeCommand)
	connection.Write([]byte(workingModeCommand))
}

func getDeviceFromIMEI(imei string) XpertDeviceData {
	client := &http.Client{}
	deviceData := XpertDeviceData{}
	var basicAuth = "Basic " + "YWZhZG1pbjphZG1pbg==" //Basic YWZhZG1pbjphZG1pbg==

	//Request Device for Id
	formattedGetDeviceRoute := getDeviceByMacRoute + "MacAddress=" + imei + "&CustomerId=" + "1"
	req, err := http.NewRequest("GET", formattedGetDeviceRoute, nil)
	req.Header.Add("Authorization", basicAuth)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	//Decode and fill Device Data
	if err := json.NewDecoder(resp.Body).Decode(&deviceData); err != nil {
		log.Fatal("Error decoding json" + err.Error())
	}
	//fmt.Println(deviceData)
	return deviceData
}

func sendSatisfiedConfigToAPI(device XpertDeviceData) {

	fmt.Println("ID:", device.Id, "PCI:", device.PendingConfigId)
	if device.PendingConfigId == 0 {
		return
	}

	client := &http.Client{}
	var basicAuth = "Basic " + "YWZhZG1pbjphZG1pbg==" //Basic YWZhZG1pbjphZG1pbg==

	deviceIdArray := `[` + fmt.Sprint(device.Id) + `]`
	//fmt.Println(deviceIdArray)
	jsonBody := []byte(deviceIdArray)
	bodyReader := bytes.NewReader(jsonBody)

	//Request Device for Id
	formattedSetConfigRoute := setConfigRoute + "ConfigId=" + fmt.Sprint(device.PendingConfigId) + "&PendingConfigID=0&CustomerId=2047"
	req, err := http.NewRequest("POST", formattedSetConfigRoute, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", basicAuth)
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer res.Body.Close()
	//resBody, err := ioutil.ReadAll(res.Body)
	//fmt.Println(string(resBody))

}
