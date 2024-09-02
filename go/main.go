package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/gousb"
	"github.com/seer-robotics/escpos"
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

	w := bufio.NewWriter(epOut)
	p := escpos.New(w)

	p.Init()

	p.SetFontSize(2, 3)
	p.SetFont("A")
	p.Write("test1")
	p.SetFont("B")

	p.WriteWEU("Привіт")

	p.SetEmphasize(1)
	p.Write("hello")
	p.Formfeed()

	p.SetUnderline(1)
	p.SetFontSize(4, 4)
	p.Write("hello")

	p.SetReverse(1)
	p.SetFontSize(2, 4)
	p.Write("hello")
	p.FormfeedN(10)

	p.SetAlign("center")
	p.Write("test")
	p.Linefeed()
	p.Write("test")
	p.Linefeed()
	p.Write("test")
	p.FormfeedD(200)

	p.Linefeed()
	p.Linefeed()
	p.Linefeed()

	p.Cut()

	w.Flush()

	fmt.Println("Print job sent successfully!")
}

func unbindKernelDriver(bus, address int) error {
	// This function attempts to unbind the kernel driver
	// Note: This requires root privileges and may not work on all systems
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo -n '%d-%d' > /sys/bus/usb/drivers/usb/unbind", bus, address))
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
