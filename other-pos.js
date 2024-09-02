let EscPosEncoder = require('esc-pos-encoder');

let encoder = new EscPosEncoder({
  width: 34,
});

let result = encoder
  .codepage('auto')
  .text('Iñtërnâtiônàlizætiøn')
  .text('διεθνοποίηση')
  .text('интернационализация')
  .encode()

console.log(result);
