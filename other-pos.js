let EscPosEncoder = require('esc-pos-encoder');

let encoder = new EscPosEncoder({
  width: 34,
});

const encodings = [
  "cp437", "cp720", "cp737", "cp775", "cp850", "cp851", "cp852", "cp853", "cp855", "cp857",
  "cp858", "cp860", "cp861", "cp862", "cp863", "cp864", "cp865", "cp866", "cp869", "cp874",
  "cp922", "cp1098", "cp1118", "cp1119", "cp1125", "cp2001", "cp3001", "cp3002", "cp3011",
  "cp3012", "cp3021", "cp3041", "cp3840", "cp3841", "cp3843", "cp3844", "cp3845", "cp3846",
  "cp3847", "cp3848", "iso885915", "iso88592", "iso88597", "rk1048", "windows1250", "windows1251",
  "windows1252", "windows1253", "windows1254", "windows1255", "windows1256", "windows1257",
  "windows1258"
];

encodings.forEach((encoding) => {
  try {
    let result = encoder
      .codepage(encoding)
      .text('інтернаціоналізація')
      .encode()

    console.log(result);
  } catch {
    // console.log('Error encoding:', encoding);
  }
});
