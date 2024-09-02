package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/gousb"
	"github.com/hennedo/escpos"
)

func main() {
	// Create a new USB context
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Find the specific device
	dev, err := ctx.OpenDeviceWithVIDPID(0x0483, 0x070b)
	if err != nil {
		log.Fatalf("Could not open device: %v", err)
	}
	defer dev.Close()

	// Attempt to detach the kernel driver
	if err := dev.SetAutoDetach(true); err != nil {
		log.Printf("Warning: Failed to set auto detach: %v", err)
	}

	// Claim the default interface (usually interface 0)
	intf, done, err := dev.DefaultInterface()
	if err != nil {
		log.Printf("Error claiming interface: %v", err)
		log.Println("Attempting to unbind kernel driver...")

		// Attempt to unbind the kernel driver
		if err := unbindKernelDriver(dev.Desc.Bus, dev.Desc.Address); err != nil {
			log.Fatalf("Failed to unbind kernel driver: %v", err)
		}

		// Try claiming the interface again
		intf, done, err = dev.DefaultInterface()
		if err != nil {
			log.Fatalf("Still unable to claim interface after unbinding: %v", err)
		}
	}
	defer done()

	// Find the OUT endpoint
	epOut, err := intf.OutEndpoint(0x01) // Assuming the first OUT endpoint
	if err != nil {
		log.Fatalf("Error finding OUT endpoint: %v", err)
	}

	p := escpos.New(epOut)

	p.Bold(true).Size(2, 2).Write("Hello World")
	p.LineFeed()
	p.Bold(false).Underline(2).Justify(escpos.JustifyCenter).Write("this is underlined")
	p.LineFeed()
	p.QRCode("https://github.com/hennedo/escpos", true, 10, escpos.QRCodeErrorCorrectionLevelH)

	// You need to use either p.Print() or p.PrintAndCut() at the end to send the data to the printer.
	p.Print()

	fmt.Println("Print job sent successfully!")
}

func unbindKernelDriver(bus, address int) error {
	// This function attempts to unbind the kernel driver
	// Note: This requires root privileges and may not work on all systems
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo -n '%d-%d' > /sys/bus/usb/drivers/usb/unbind", bus, address))
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
