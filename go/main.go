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

var encodings = []string{
	"PC437", "Katakana", "PC850", "PC860", "PC863", "PC865", "CP1252",
	"ISO-8859-1", "ISO-8859-7", "ISO-8859-8", "ISO-8859-2", "Windows-1256",
	"Windows-1252", "PC866", "PC852", "PC858", "Windows-1256", "ISO-8859-13",
	"Windows-1256", "PT151",
}

func WriteToPrinter(e *escpos.Escpos, data string, encoding string) (int, error) {
	cd, err := iconv.Open(encoding, "utf-8")
	if err != nil {
		return 0, fmt.Errorf("error opening iconv: %v", err)
	}
	defer cd.Close()
	weu := cd.ConvString(data)
	return e.WriteRaw([]byte(weu))
}

func main() {
	ctx := gousb.NewContext()
	defer ctx.Close()

	dev, err := ctx.OpenDeviceWithVIDPID(0x0483, 0x070b)
	if err != nil {
		log.Fatalf("Could not open device: %v", err)
	}
	defer dev.Close()

	if err := dev.SetAutoDetach(true); err != nil {
		log.Printf("Warning: Failed to set auto detach: %v", err)
	}

	intf, done, err := dev.DefaultInterface()
	if err != nil {
		log.Printf("Error claiming interface: %v", err)
		log.Println("Attempting to unbind kernel driver...")
		if err := unbindKernelDriver(dev.Desc.Bus, dev.Desc.Address); err != nil {
			log.Fatalf("Failed to unbind kernel driver: %v", err)
		}
		intf, done, err = dev.DefaultInterface()
		if err != nil {
			log.Fatalf("Still unable to claim interface after unbinding: %v", err)
		}
	}
	defer done()

	epOut, err := intf.OutEndpoint(0x01)
	if err != nil {
		log.Fatalf("Error finding OUT endpoint: %v", err)
	}

	w := bufio.NewWriter(epOut)
	p := escpos.New(w)

	p.Init()
	p.SetFontSize(1, 1)
	p.SetAlign("center")
	p.SetFont("A")

	for _, encoding := range encodings {
		p.WriteRaw([]byte(fmt.Sprintf("Encoding: %s\n", encoding)))
		_, err := WriteToPrinter(p, "Привіт", encoding)
		if err != nil {
			p.WriteRaw([]byte(fmt.Sprintf("Error: %v\n", err)))
		}
		p.Linefeed()
		p.Linefeed()
	}

	p.FormfeedN(5)
	p.Cut()

	w.Flush()

	fmt.Println("Print job sent successfully!")
}

func unbindKernelDriver(bus, address int) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo -n '%d-%d' > /sys/bus/usb/drivers/usb/unbind", bus, address))
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
