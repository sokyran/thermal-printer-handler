package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/gousb"
)

// ESC/POS commands
var (
	ESC  = []byte{0x1B}
	GS   = []byte{0x1D}
	INIT = []byte{0x1B, 0x40}       // Initialize printer
	CUT  = []byte{0x1D, 0x56, 0x41} // Cut paper
	LF   = []byte{0x0A}             // Line feed
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

	// Send initialization command
	_, err = epOut.Write(INIT)
	if err != nil {
		log.Fatalf("Error initializing printer: %v", err)
	}

	// Print some text
	text := []byte("Hello, Thermal Printer!\n")
	_, err = epOut.Write(text)
	if err != nil {
		log.Fatalf("Error writing text: %v", err)
	}

	// Feed a line
	_, err = epOut.Write(LF)
	if err != nil {
		log.Fatalf("Error feeding line: %v", err)
	}

	// Cut the paper
	_, err = epOut.Write(CUT)
	if err != nil {
		log.Fatalf("Error cutting paper: %v", err)
	}

	fmt.Println("Print job sent successfully!")
}

func unbindKernelDriver(bus, address int) error {
	// This function attempts to unbind the kernel driver
	// Note: This requires root privileges and may not work on all systems
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo -n '%d-%d' > /sys/bus/usb/drivers/usb/unbind", bus, address))
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
