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
	ESC       = []byte{0x1B}
	GS        = []byte{0x1D}
	INIT      = []byte{0x1B, 0x40}       // Initialize printer
	CUT       = []byte{0x1D, 0x56, 0x41} // Cut paper
	LF        = []byte{0x0A}             // Line feed
	nodeArray = []byte{
		27, 116, 0, 73, 164, 116, 137, 114, 110, 131, 116,
		105, 147, 110, 133, 108, 105, 122, 145, 116, 105, 27,
		116, 19, 155, 110, 27, 116, 0, 235, 27, 116, 38,
		227, 222, 226, 231, 233, 234, 233, 159, 225, 236, 225,
		27, 116, 34, 183, 212, 229, 168, 225, 212, 160, 164,
		183, 214, 212, 160, 208, 183, 243, 160, 164, 183, 222,
		10, 13,
	}
)

var arr1 = []byte{
	27, 116, 0, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr2 = []byte{
	27, 116, 32, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr3 = []byte{
	27, 116, 14, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr4 = []byte{
	27, 116, 33, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr5 = []byte{
	27, 116, 2, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr6 = []byte{
	27, 116, 11, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr7 = []byte{
	27, 116, 18, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr8 = []byte{
	27, 116, 12, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr9 = []byte{
	27, 116, 34, 138, 212, 229, 168, 225, 212, 160, 164, 138, 214, 212, 160, 208, 138, 243, 160, 164, 138, 222,
}

var arr10 = []byte{
	27, 116, 13, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr11 = []byte{
	27, 116, 19, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr12 = []byte{
	27, 116, 3, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr13 = []byte{
	27, 116, 35, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr14 = []byte{
	27, 116, 36, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr15 = []byte{
	27, 116, 4, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr16 = []byte{
	27, 116, 37, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr17 = []byte{
	27, 116, 5, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr18 = []byte{
	27, 116, 17, 63, 173, 226, 165, 224, 173, 160, 230, 63, 174, 173, 160, 171, 63, 167, 160, 230, 63, 239,
}

var arr19 = []byte{
	27, 116, 38, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr20 = []byte{
	27, 116, 41, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr21 = []byte{
	27, 116, 42, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr22 = []byte{
	27, 116, 43, 63, 173, 226, 165, 224, 173, 160, 230, 63, 174, 173, 160, 171, 63, 167, 160, 230, 63, 239,
}

var arr23 = []byte{
	27, 116, 44, 247, 173, 226, 165, 224, 173, 160, 230, 247, 174, 173, 160, 171, 247, 167, 160, 230, 247, 239,
}

var arr24 = []byte{
	27, 116, 40, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr25 = []byte{
	27, 116, 39, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr26 = []byte{
	27, 116, 15, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr27 = []byte{
	27, 116, 53, 179, 237, 242, 229, 240, 237, 224, 246, 179, 238, 237, 224, 235, 179, 231, 224, 246, 179, 255,
}

var arr28 = []byte{
	27, 116, 45, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr29 = []byte{
	27, 116, 46, 179, 237, 242, 229, 240, 237, 224, 246, 179, 238, 237, 224, 235, 179, 231, 224, 246, 179, 255,
}

var arr30 = []byte{
	27, 116, 16, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr31 = []byte{
	27, 116, 47, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr32 = []byte{
	27, 116, 48, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr33 = []byte{
	27, 116, 49, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr34 = []byte{
	27, 116, 50, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr35 = []byte{
	27, 116, 51, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arr36 = []byte{
	27, 116, 52, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63,
}

var arrays = [][]byte{
	arr1, arr2, arr3, arr4, arr5, arr6, arr7, arr8, arr9, arr10,
	arr11, arr12, arr13, arr14, arr15, arr16, arr17, arr18, arr19, arr20,
	arr21, arr22, arr23, arr24, arr25, arr26, arr27, arr28, arr29, arr30,
	arr31, arr32, arr33, arr34, arr35, arr36,
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

	for _, arr := range arrays {
		_, err = epOut.Write(arr)
		if err != nil {
			log.Fatalf("Error writing text: %v", err)
		}
		epOut.Write(LF)
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
