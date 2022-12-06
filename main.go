package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	host           = "172.31.51.82"
	port           = "8000"
	connectionType = "tcp"
	locationRoute  = "http://3.212.201.170:802/XpertRestApi/api/location_data" //http://52.45.17.177:802/XpertRestApi/api/location_data //http://3.212.201.170:802/XpertRestApi/api/location_data
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

		fmt.Println(message)

		// Sorting
		if strings.Contains(message, "AP00") { // AP00 = Connection
			deviceIMEI = getIMEIFromAP00(message)
			go connection.Write([]byte("IWBP00," + dtg + ",4#"))

		} else if strings.Contains(message, "AP03") { // AP03 = Heartbeat
			go connection.Write([]byte("IWBP03#"))

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
	fmt.Printf("XpertRestAPI response: %s\n", resBody)
}
