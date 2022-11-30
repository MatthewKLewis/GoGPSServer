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
	locationRoute  = "http://3.212.201.170:802/XpertRestApi/api/location_data"
	alertRoute     = "http://3.212.201.170:802/XpertRestApi/api/alert_data"
)

func main() {
	socketServer, err := net.Listen(connectionType, host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
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

	for {
		buffer := make([]byte, 1024)
		dtg := time.Now().Format("01/30/2006 15:04:05")

		// Waits here for next message...
		mLen, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			connection.Close()
			return
		}

		// Handles Message
		message := string(buffer[:mLen])
		//fmt.Println(message)

		// Sorting
		if strings.Contains(message, "AP00") { // AP00 = Connection
			deviceIMEI = getIMEIFromAP00(message)
			connection.Write([]byte("IWBP00," + dtg + ",4#"))

		} else if strings.Contains(message, "AP01") { // AP01 = Location
			var packetData, err = getJSONFromAP01(message, deviceIMEI)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			body, err := json.Marshal(packetData)
			if err != nil {
				fmt.Println("Error marshaling json:", err.Error())
				return
			}

			res, err := http.Post(alertRoute, "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			resBody, _ := ioutil.ReadAll(res.Body)
			fmt.Printf("client: response body: %s\n", resBody)

		} else if strings.Contains(message, "AP03") { // AP03 = Heartbeat
			connection.Write([]byte("IWBP03#"))

		} else if strings.Contains(message, "AP10") { // AP10 = Alert
			var packetData, err = getJSONFromAP10(message, deviceIMEI)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			_, err = json.Marshal(packetData)
			if err != nil {
				fmt.Println("Error marshaling json:", err.Error())
				return
			}

			fmt.Println(packetData)
			// http.Post(alertRoute, "application/json", bytes.NewBuffer(body))

		} else if strings.Contains(message, "LK") { // LK = Link
			fmt.Println("-- LK")
			deviceIMEI = getIMEIFromLK(message)
			stringToSend := "[3G*" + deviceIMEI + "*0002*LK]"
			fmt.Println(stringToSend)
			connection.Write([]byte(stringToSend))

		} else if strings.Contains(message, "CUSTOMER") { // CUSTOMER = Location for LK-type messages
			fmt.Println("-- CUSTOMER...")

		} else {
			//fmt.Println("-- Other message...")
		}
	}
}
