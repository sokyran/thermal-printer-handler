const escpos = require('escpos');
// install escpos-usb adapter module manually
escpos.USB = require('escpos-usb');
// Select the adapter based on your printer type
const device = new escpos.USB();
// const device  = new escpos.Network('localhost');
// const device  = new escpos.Serial('/dev/usb/lp0');

const options = { encoding: "windows-1251" }

const printer = new escpos.Printer(device, options);

const multilineText = `
— А ви заблокували Юру ? Він оце без вас поїхав
— Тоді ще Миру треба…
— Давайте!

Влад прийшов на коворкинг до Насти та Єви
`;


device.open(function (error) { printer.align('ct').encode('ibm857').text('The quick brown fox jumps over the lazy dog').text('Ç Ğ I İ Ö Ş Ü ü ş ö i ı ğ ç Â â î Î Û û').qrimage('https://github.com/song940/node-escpos', function (err) { this.cut(); this.close(); }); });

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
