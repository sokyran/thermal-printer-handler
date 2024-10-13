import { findByIds } from 'usb';
import ReceiptPrinterEncoder from '@point-of-sale/receipt-printer-encoder';
import getPixels from 'get-pixels';

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

let encoder = new ReceiptPrinterEncoder({
  width: 32,
});

const imagePath = './svyato.png';

let pixels = await new Promise(resolve => {
  getPixels(imagePath, (err, pixels) => {
    resolve(pixels);
  });
});

let result = encoder
  .initialize()
  .text('')
  .newline()
  .image(pixels, 400, 400, 'atkinson')
  .newline()
  .newline()
  .newline()
  .newline()
  .newline()
  .text('')
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
