const escpos = require("escpos");
// install escpos-usb adapter module manually
escpos.USB = require("escpos-usb");
// Select the adapter based on your printer type
const device = new escpos.USB();

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
