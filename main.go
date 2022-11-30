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
	fmt.Println("Received: ", string(buffer[:mLen]), " from: ", connection.LocalAddr().String())

	message := string(buffer[:mLen])

	if strings.Contains(message, "AP00") { // AP00 = Connection
		fmt.Print(parseAP00(string(buffer[:mLen])))

	} else if strings.Contains(message, "AP01") { // AP01 = Location?
		fmt.Print(parseAP01(string(buffer[:mLen])))

	} else if strings.Contains(message, "AP10") { // AP10 = Alert?
		fmt.Print(parseAP10(string(buffer[:mLen])))

	}
	//connection.Close()
}
