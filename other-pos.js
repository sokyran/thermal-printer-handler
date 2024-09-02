import EscPosEncoder from 'esc-pos-encoder';
import { getDeviceList, findByIds } from 'usb';

const encodings = [
  "cp437", "cp720", "cp737", "cp775", "cp850", "cp851", "cp852", "cp853", "cp855", "cp857",
  "cp858", "cp860", "cp861", "cp862", "cp863", "cp864", "cp865", "cp866", "cp869", "cp874",
  "cp922", "cp1098", "cp1118", "cp1119", "cp1125", "cp2001", "cp3001", "cp3002", "cp3011",
  "cp3012", "cp3021", "cp3041", "cp3840", "cp3841", "cp3843", "cp3844", "cp3845", "cp3846",
  "cp3847", "cp3848", "iso885915", "iso88592", "iso88597", "rk1048", "windows1250", "windows1251",
  "windows1252", "windows1253", "windows1254", "windows1255", "windows1256", "windows1257",
  "windows1258"
];

// encodings.forEach((encoding) => {
//   try {
//     let result = encoder
//       .codepage(encoding)
//       .text('інтернаціоналізація')
//       .encode()

//     console.log(result);
//   } catch {
//     // console.log('Error encoding:', encoding);
//   }
// });

// Step 1: Find the Printer
const printerVendorId = 1155; // Replace with your printer's vendor ID
const printerProductId = 1803; // Replace with your printer's product ID

const device = findByIds(printerVendorId, printerProductId);

if (!device) {
  console.error('Printer not found');
  process.exit(1);
}

// Step 2: Open the Printer Connection
device.open();

const intr = device.interfaces[0];
if (intr.isKernelDriverActive()) {
  intr.detachKernelDriver();
}

intr.claim();

const endpoint = intr.endpoints[0];

if (!endpoint || !endpoint.direction === 'out') {
  console.error('Printer endpoint not found');
  process.exit(1);
}

// Step 3: Create and Encode the Print Command
const encoder = new EscPosEncoder({
  width: 34,
});
const result = encoder
  .initialize()
  .text('Hello World')
  .newline()
  .encode();

// Step 4: Send the Data to the Printer
endpoint.transfer(result, (error) => {
  if (error) {
    console.error('Failed to print:', error);
  } else {
    console.log('Printed successfully');
  }

  // Step 5: Release the Interface and Close the Device
  intr.release(true, () => {
    device.close();
  });
});
