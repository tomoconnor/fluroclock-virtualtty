package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"go.bug.st/serial"
)

func sendResponse(message string, port *serial.Port) {
	panel_id := os.Getenv("PANEL_ID")
	if panel_id == "" {
		log.Println("PANEL_ID not set, defaulting to 1")
		panel_id = "1"
	}
	panel_id = strings.TrimSpace(panel_id)

	fmt.Print("Message: ")
	fmt.Println(message)
	// fmt.Println("\nEOF")
	response := "+" + panel_id + "-"

	sp := *port
	if strings.Contains(message, "+?-") {
		m, err := sp.Write([]byte(response))
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Sent:", m, "bytes")
	}

}

func main() {

	// Open the serial port
	serial_port_name := os.Getenv("SERIAL_PORT_NAME")
	if serial_port_name == "" {
		log.Fatal("SERIAL_PORT_NAME not set")
	}

	baud := os.Getenv("BAUD")
	if baud == "" {
		log.Fatal("BAUD not set")
	}

	ibaud, err := strconv.Atoi(baud)
	if err != nil {
		log.Fatal("BAUD must be an integer")
	}

	mode := serial.Mode{
		BaudRate: int(ibaud),
	}
	port, err := serial.Open(serial_port_name, &mode)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()
	message := make([]byte, 1)

	buff := make([]byte, 100)
	for {
		// Reads up to 100 bytes
		n, err := port.Read(buff)
		// fmt.Println("Got", n, "bytes")
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		message = append(message, buff[:n]...)
		// fmt.Printf("%s", string(buff[:n]))
		// fmt.Printf("buffer: %s", string(buff))
		// If we receive a - stop reading
		if strings.Contains(string(buff[:n]), "-") {
			sendResponse(string(message), &port)
			message = make([]byte, 1) // Reset the message
		}
	}

}
