package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	host           = "localhost"
	port           = "8000"
	connectionType = "tcp"
	locationRoute  = "http://3.212.201.170:802/XpertRestApi/api/location_data"
	alertRoute     = "http://3.212.201.170:802/XpertRestApi/api/alert_data"
)

func main() {

	fmt.Println("Server Running...")
	server, err := net.Listen(connectionType, host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()

	fmt.Println("Listening on " + host + ":" + port)
	fmt.Println("Waiting for client...")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("client connected")
		go threadedClientConnectionHandler(connection)
	}
}

func threadedClientConnectionHandler(connection net.Conn) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Received: ", string(buffer[:mLen]))

	message := string(buffer[:mLen])

	// AP00 = Connection
	if strings.Contains(message, "AP00") {
		parseAP00(string(buffer[:mLen]))
	} else if strings.Contains(message, "AP01") {
		parseAP01(string(buffer[:mLen]))
	} else if strings.Contains(message, "AP10") {
		parseAP10(string(buffer[:mLen]))
	}

	_, err = connection.Write([]byte("Thanks! Got your message:" + string(buffer[:mLen])))
	//connection.Close()
}

func parseAP00(msg string) string {
	return msg
}

func parseAP01(msg string) string {
	return msg
}

func parseAP10(msg string) string {
	return msg
}
