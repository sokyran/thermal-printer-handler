package main

import (
	"fmt"
	"os"
)

func main() {
	// Open the USB device file
	deviceFile := "/dev/bus/usb/001/002" // Replace with the correct device path
	file, err := os.OpenFile(deviceFile, os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Error opening USB device: %v\n", err)
		return
	}
	defer file.Close()

	file.Write([]byte{0x1b, 0x40})

	// p := escpos.New(file)

	// // p.Verbose = true

	// p.Init()
	// p.SetFontSize(2, 3)
	// p.SetFont("A")
	// p.Write("test1")
	// p.SetFont("B")
	// p.Write("test2")

	// p.SetEmphasize(1)
	// p.Write("hello")
	// p.Formfeed()

	// p.SetUnderline(1)
	// p.SetFontSize(4, 4)
	// p.Write("hello")

	// p.SetReverse(1)
	// p.SetFontSize(2, 4)
	// p.Write("hello")
	// p.FormfeedN(10)

	// p.SetAlign("center")
	// p.Write("test")
	// p.Linefeed()
	// p.Write("test")
	// p.Linefeed()
	// p.Write("test")
	// p.FormfeedD(200)

	// p.Cut()
}
