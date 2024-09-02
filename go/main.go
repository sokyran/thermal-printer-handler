package main

import (
	"fmt"
	"log"

	"github.com/google/gousb"
)

func main() {
	// Create a new USB context
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Get a list of all USB devices
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		// This function is called for each device.
		// Returning true means the device should be opened.
		return true
	})

	if err != nil {
		log.Fatalf("Error opening devices: %v", err)
	}

	// Ensure all devices are closed when we're done
	defer func() {
		for _, d := range devices {
			d.Close()
		}
	}()

	// Print information about each device
	for i, dev := range devices {
		fmt.Printf("Device %d:\n", i)
		fmt.Printf("  Bus: %d\n", dev.Desc.Bus)
		fmt.Printf("  Address: %d\n", dev.Desc.Address)
		fmt.Printf("  Speed: %s\n", dev.Desc.Speed)
		fmt.Printf("  Vendor ID: %s\n", dev.Desc.Vendor)
		fmt.Printf("  Product ID: %s\n", dev.Desc.Product)

		// Get the device configuration
		cfg, err := dev.Config(1)
		if err != nil {
			log.Printf("Error getting config: %v", err)
			continue
		}

		// Print information about each interface
		for _, intf := range cfg.Desc.Interfaces {
			for _, ifSetting := range intf.AltSettings {
				fmt.Printf("    Interface %d:\n", ifSetting.Number)
				fmt.Printf("      Alt Setting: %d\n", ifSetting.Alternate)
				fmt.Printf("      Class: %d\n", ifSetting.Class)
				fmt.Printf("      SubClass: %d\n", ifSetting.SubClass)
				fmt.Printf("      Protocol: %d\n", ifSetting.Protocol)
			}
		}

		fmt.Println()
	}
}
