package main

import (
	"fmt"

	"github.com/google/gousb"
)

func main() {
	ctx := gousb.NewContext()

	defer ctx.Close()

	dev, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		fmt.Printf("%s\n", desc)
		return desc.Vendor == 0x0416 && desc.Product == 0x5011
	})

	// if err != nil {
	// 	log.Fatalf("Could not open a device: %v", err)
	// }
	// defer dev.Close()

	// // Claim the default interface using a convenience function.
	// // The default interface is always #0 alt #0 in the currently active
	// // config.
	// intf, done, err := dev.DefaultInterface()
	// if err != nil {
	// 	log.Fatalf("%s.DefaultInterface(): %v", dev, err)
	// }
	// defer done()

	// // Open an OUT endpoint.
	// ep, err := intf.OutEndpoint(7)
	// if err != nil {
	// 	log.Fatalf("%s.OutEndpoint(7): %v", intf, err)
	// }

	// // Create a new printer
	// p := escpos.New(ep)

	// p.Bold(true).Size(2, 2).Write("Hello World")
	// p.LineFeed()
	// p.Bold(false).Underline(2).Justify(escpos.JustifyCenter).Write("this is underlined")
	// p.LineFeed()
	// p.QRCode("https://github.com/hennedo/escpos", true, 10, escpos.QRCodeErrorCorrectionLevelH)

	// // You need to use either p.Print() or p.PrintAndCut() at the end to send the data to the printer.
	// p.Print()
}
