package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/gousb"
	"github.com/qiniu/iconv"
	"github.com/seer-robotics/escpos"
)

var encodings = []string{
	"ASCII", "UTF-8", "ISO-8859-1", "ISO-8859-2", "ISO-8859-3", "ISO-8859-4", "ISO-8859-5",
	"ISO-8859-6", "ISO-8859-7", "ISO-8859-8", "ISO-8859-9", "ISO-8859-10", "ISO-8859-11",
	"ISO-8859-13", "ISO-8859-14", "ISO-8859-15", "ISO-8859-16", "KOI8-R", "KOI8-U", "KOI8-RU",
	"CP1250", "CP1251", "CP1252", "CP1253", "CP1254", "CP1255", "CP1256", "CP1257", "CP1258",
	"CP850", "CP862", "CP866", "MACINTOSH", "MACCENTRALEUROPE", "MACICELAND", "MACCROATIAN",
	"MACROMANIA", "MACCYRILLIC", "MACUKRAINE", "MACGREEK", "MACTURKISH", "MACHEBREW", "MACARABIC",
	"MACTHAI", "HP-ROMAN8", "NEXTSTEP", "ARMSCII-8", "GEORGIAN-ACADEMY", "GEORGIAN-PS", "KOI8-T",
	"PT154", "RK1048", "MULELAO-1", "CP1133", "TIS-620", "CP874", "VISCII", "TCVN", "SHIFT-JIS",
	"CP932", "ISO-2022-JP", "EUC-CN", "GBK", "CP936", "GB18030", "ISO-2022-CN", "HZ", "EUC-TW",
	"BIG-5", "CP950", "BIG5-HKSCS", "EUC-KR", "CP949", "CP1361", "ISO-2022-KR", "CP856", "CP922",
	"CP943", "CP1046", "CP1124", "CP1129", "CP1161", "CP1162", "CP1163", "DEC-KANJI", "DEC-HANYU",
	"CP437", "CP737", "CP775", "CP852", "CP853", "CP855", "CP857", "CP858", "CP860", "CP861",
	"CP863", "CP864", "CP865", "CP869", "CP1125", "EUC-JISX0213", "SHIFT_JISX0213", "ISO-2022-JP-3",
	"BIG5-2003", "TDS565", "ATARI", "RISCOS-LATIN1",
}

func WriteToPrinter(e *escpos.Escpos, data string, encoding string) (int, error) {
	cd, err := iconv.Open(encoding, "UTF-8")
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
	p.FormfeedN(3)

	for _, encoding := range encodings {
		p.WriteRaw([]byte(fmt.Sprintf("Encoding: %s\n", encoding)))
		_, err := WriteToPrinter(p, "Привіт", encoding)
		if err != nil {
			if strings.Contains(err.Error(), "iconv_open") {
				p.WriteRaw([]byte("Encoding not supported by iconv\n"))
			} else {
				p.WriteRaw([]byte(fmt.Sprintf("Error: %v\n", err)))
			}
		}
		p.Linefeed()
	}

	p.Cut()

	w.Flush()

	fmt.Println("Print job sent successfully!")
}

func unbindKernelDriver(bus, address int) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo -n '%d-%d' > /sys/bus/usb/drivers/usb/unbind", bus, address))
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
