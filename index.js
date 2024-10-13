import { USB, Printer } from "escpos";
// install escpos-usb adapter module manually
import * as USB from 'escpos-usb';
// Select the adapter based on your printer type
const device = new USB();

const options = { encoding: "windows-1251" };

const printer = new Printer(device, options);

device.open(function (error) {
  printer
    .align("ct")
    .image("./image.jpg")
    .cut()
    .close();
});

// device.open(function (error) {
//   printer
//     .encode('cp866')
//     .font('a')
//     .align('LT')
//     .style('bu')
//     .size(1, 1)

//   multilineText.split('\n').forEach((line) => {
//     printer.text(line);
//   });

//   printer.cut().close();
// });
