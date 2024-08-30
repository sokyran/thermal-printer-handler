package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hennedo/escpos"
)

func main() {
	// Open the USB device file
	deviceFile := "/dev/bus/usb/001/002"
	file, err := os.OpenFile(deviceFile, os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Error opening USB device: %v\n", err)
		return
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	// Create a new escpos printer using the USB device
	p := escpos.New(w) // `file` is an `io.Writer`

	p.Bold(true).Size(2, 2).Write("Hello World")
	p.LineFeed()
	p.Bold(false).Underline(2).Justify(escpos.JustifyCenter).Write("this is underlined")
	p.LineFeed()
	p.QRCode("https://github.com/hennedo/escpos", true, 10, escpos.QRCodeErrorCorrectionLevelH)

	// You need to use either p.Print() or p.PrintAndCut() at the end to send the data to the printer.
	p.PrintAndCut()
}
