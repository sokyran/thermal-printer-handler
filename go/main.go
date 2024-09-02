package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/gousb"
	"github.com/qiniu/iconv"
	"github.com/seer-robotics/escpos"
)

func WriteToPrinter(e *escpos.Escpos, data string) (int, error) {
	cd, err := iconv.Open("utf-8", "utf-8")
	if err != nil {
		log.Fatal(err)
	}
	defer cd.Close()
	weu := cd.ConvString(data)
	return e.WriteRaw([]byte(weu))
}

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

	p.SetFontSize(1, 1)
	p.SetAlign("center")
	p.SetFont("A")
	p.Write("test1")
	p.Linefeed()

	WriteToPrinter(p, "test2")
	WriteToPrinter(p, "ПРИВІТ")

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
