const escpos = require("escpos");
// install escpos-usb adapter module manually
const { findByIds } = require("usb");
// Select the adapter based on your printer type
// Step 1: Find the Printer
const printerVendorId = 1155; // Replace with your printer's vendor ID
const printerProductId = 1803; // Replace with your printer's product ID

const device = findByIds(printerVendorId, printerProductId);

const options = { encoding: "windows-1251" };

const printer = new escpos.Printer(device, options);

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
