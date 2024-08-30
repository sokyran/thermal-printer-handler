const escpos = require('escpos');
// install escpos-usb adapter module manually
escpos.USB = require('escpos-usb');
// Select the adapter based on your printer type
const device = new escpos.USB();
// const device  = new escpos.Network('localhost');
// const device  = new escpos.Serial('/dev/usb/lp0');

const options = { encoding: "utf-8" }
// encoding is optional

const printer = new escpos.Printer(device, options);

const multilineText = `
— А ви заблокували Юру ? Він оце без вас поїхав
— Тоді ще Миру треба…
— Давайте!

Влад прийшов на коворкинг до Насти та Єви
`;


device.open(function (error) {
  printer
    .font('a')
    .align('ct')
    .style('bu')
    .size(1, 1)

  multilineText.split('\n').forEach((line) => {
    printer.text(line);
  });

  printer.cut().close();
});
