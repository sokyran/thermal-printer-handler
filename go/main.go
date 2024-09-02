package main

import (
	"fmt"
	"os"

	"github.com/hennedo/escpos"
	"go.bug.st/serial"
)

func main() {
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open("/dev/bus/usb/001/002", mode)

	if err != nil {
		fmt.Print("Error opening serial port: ", err)
		return
	}

	fmt.Print(port, err)

	// Open the USB device file
	deviceFile := "/dev/bus/usb/001/002" // Replace with the correct device path
	file, err := os.OpenFile(deviceFile, os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Error opening USB device: %v\n", err)
		return
	}
	defer file.Close()

	// Create a new printer
	p := escpos.New(file)

	p.Bold(true).Size(2, 2).Write("Hello World")
	p.LineFeed()
	p.Bold(false).Underline(2).Justify(escpos.JustifyCenter).Write("this is underlined")
	p.LineFeed()
	p.QRCode("https://github.com/hennedo/escpos", true, 10, escpos.QRCodeErrorCorrectionLevelH)

	// You need to use either p.Print() or p.PrintAndCut() at the end to send the data to the printer.
	p.Print()
}
