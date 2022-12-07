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
	locationRoute         = "http://3.212.201.170:802/XpertRestApi/api/location_data" //http://52.45.17.177:802/XpertRestApi/api/location_data //http://3.212.201.170:802/XpertRestApi/api/location_data
	getDeviceByMacRoute   = "http://3.212.201.170:802/XpertRestApi/api/Device/GetByMacAddress?"
	getConfigPendingRoute = "http://3.212.201.170:802/XpertRestApi/api/Device/GetConfigPendingByDeviceId?"
	setConfigRoute        = "http://3.212.201.170:802/XpertRestApi/api/Device/SetDeviceConfigurations?  CustomerId=2047 &ConfigID=1235 &PendingConfigID=0"
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

		//if there is an IMEI, call to API for GetConfigByMAC / GET
		//get config by mac
		//set config as active

		//fmt.Println(message)

		if deviceIMEI != "" {
			go getConfigurationForTag(deviceIMEI)
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
			go sendToAPI(packetData)

		} else if strings.Contains(message, "AP01") { // AP01 = Location
			packetData, err = getJSONFromAP01("IW"+message, deviceIMEI)
			if err != nil {
				fmt.Println("Breaking from While Loop, No Lat Lon!")
				break
			}
			go sendToAPI(packetData)

		} else if strings.Contains(message, "AP10") { // AP10 = Alert
			packetData, err = getJSONFromAP10(message, deviceIMEI)
			handleError(err)
			go connection.Write([]byte("IWBP10#"))
			go sendToAPI(packetData)

		} else if strings.Contains(message, "LK") { // LK = Link
			// deviceIMEI = getIMEIFromLK(message)
			// [3G*8800000015*0009*UPLOAD,600] [3G*8800000015*0027*SOS,00000000000,00000000000,00000000000]
			// connection.Write([]byte("[3G*" + deviceIMEI + "*0002*LK]"))
			// connection.Write([]byte("[3G*" + deviceIMEI + "*0009*UPLOAD,10]"))
			// connection.Write([]byte("[3G*" + deviceIMEI + "*0027*SOS,14438137623,00000000000,00000000000]"))
		} else if strings.Contains(message, "ALCUSTOMER") { // ALCUSTOMER = Alarm
			// packetData, err = getJSONFromCUSTOMER(message, deviceIMEI)
			// handleError(err)
			// sendToAPI(packetData)
		} else if strings.Contains(message, "UDCUSTOMER") { // UDCUSTOMER = Location
			// packetData, err = getJSONFromCUSTOMER(message, deviceIMEI)
			// handleError(err)
			// sendToAPI(packetData)
		}
	}
}

func sendToAPI(packet PreciseGPSData) {
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
func getConfigurationForTag(imei string) {

	client := &http.Client{}
	deviceData := XpertDeviceData{}
	configData := XpertConfigData{}
	var basicAuth = "Basic " + "YWZhZG1pbjphZG1pbg==" //Basic YWZhZG1pbjphZG1pbg==

	//Request
	formattedGetDeviceRoute := getDeviceByMacRoute + "MacAddress=" + imei + "&CustomerId=" + "2047"
	req, err := http.NewRequest("GET", formattedGetDeviceRoute, nil)
	req.Header.Add("Authorization", basicAuth)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	//Decode
	if err := json.NewDecoder(resp.Body).Decode(&deviceData); err != nil {
		log.Fatal("Error decoding json" + err.Error())
	}
	fmt.Println("Device ID", deviceData.Id)

	//Request
	formattedGetConfigRoute := getConfigPendingRoute + "DeviceId=" + fmt.Sprint(deviceData.Id) + "&CustomerId=" + "2047"
	fmt.Println(formattedGetConfigRoute)
	req, err = http.NewRequest("GET", formattedGetConfigRoute, nil)
	req.Header.Add("Authorization", basicAuth)
	resp, err = client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	//Decode
	if err := json.NewDecoder(resp.Body).Decode(&configData); err != nil {
		log.Fatal("Error decoding json" + err.Error())
	}
	fmt.Println("Config ID", configData.Id)

	// res, err = http.Get(formattedGetConfigRoute + "DeviceId=" + imei + "&CustomerId=" + "2047")
	// handleError(err)
	// resBody, _ = ioutil.ReadAll(res.Body)
	// fmt.Printf("Config Update API response: %s\n", resBody)

}
