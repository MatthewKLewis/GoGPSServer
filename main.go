package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	host           = "172.31.51.82" // 127.0.0.1 // 172.31.51.82
	port           = "8000"
	connectionType = "tcp"
	locationRoute  = "http://3.212.201.170:802/XpertRestApi/api/location_data"
	alertRoute     = "http://3.212.201.170:802/XpertRestApi/api/alert_data"
)

var numberOfConnections = 0

func main() {
	connectionMap := make(map[string]int)

	fmt.Println("Server Running...")
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
		numberOfConnections++
		connectionMap[connection.RemoteAddr().String()] = numberOfConnections
		go threadedClientConnectionHandler(connection)
	}
}

func threadedClientConnectionHandler(connection net.Conn) {
	for {
		buffer := make([]byte, 1024)
		dtg := time.Now().Format("01/02/2006 15:04:05")

		// Waits here for next message...
		mLen, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			connection.Close()
			return
		}

		// Handles Message
		fmt.Println(dtg, "Received from:", connection.RemoteAddr().String())
		message := string(buffer[:mLen])
		if strings.Contains(message, "AP00") { // AP00 = Connection
			fmt.Println("-- AP00")
			fmt.Println("-- IMEI:", getIMEIFromAP00(string(buffer[:mLen])))

			connection.Write([]byte("IWBP00," + dtg + ",4#"))

		} else if strings.Contains(message, "AP01") { // AP01 = Location?
			fmt.Println("-- AP01")

			var packetData = getJSONFromAP01(string(buffer[:mLen]))
			body, err := json.Marshal(packetData)
			if err != nil {
				fmt.Println("Error marshaling json:", err.Error())
				return
			}
			http.Post(locationRoute, "application/json", bytes.NewBuffer(body))

		} else if strings.Contains(message, "AP03") { // AP03 = Heartbeat
			fmt.Println("-- AP03")
			fmt.Println("-- Heartbeat")

			connection.Write([]byte("IWBP03#"))

		} else if strings.Contains(message, "AP10") { // AP10 = Alert?
			fmt.Println("-- AP10")
			fmt.Println("--", getJSONFromAP10(string(buffer[:mLen])))

			//stringToSend := "Hello"
			//connection.Write([]byte(stringToSend))

		} else if strings.Contains(message, "LK") { // LK = Link?
			fmt.Println("-- LK")
			fmt.Println("-- IMEI:", getIMEIFromLK(string(buffer[:mLen])))

			//stringToSend := "Hello"
			//connection.Write([]byte(stringToSend))

		} else if strings.Contains(message, "CUSTOMER") { // LK = Link?
			fmt.Println("-- CUSTOMER...")

			//stringToSend := "Hello"
			//connection.Write([]byte(stringToSend))

		} else {
			fmt.Println("-- Other message...")
		}
	}
}
