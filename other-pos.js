import { findByIds } from 'usb';
import ReceiptPrinterEncoder from '@point-of-sale/receipt-printer-encoder';
import getPixels from 'get-pixels';
import sizeOf from 'image-size';
import fs from 'fs';
import path from 'path';

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

const encoder = new ReceiptPrinterEncoder({
  width: 32,
});

let result = encoder.initialize();


// Directory where the images are located
const imagesDir = './images/';

const counter = 0;

const files = fs.readdirSync(imagesDir);

console.log('Printing images:', files);

const imagesPaths = files.map(file => path.join(imagesDir, file));

for (let imagePath of [imagesPaths[counter]]) {
  const dimensions = sizeOf(imagePath);

  const aspectRatio = dimensions.width / dimensions.height;
  const newHeight = 400 / aspectRatio;
  const roundedHeight = Math.round(newHeight / 8) * 8;

  let pixels = await new Promise(resolve => {
    getPixels(imagePath, (err, pixels) => {
      resolve(pixels);
    });
  });

  result
    .text('')
    .newline()
    .image(pixels, 400, roundedHeight, 'atkinson')
    .newline()
    .text('')
}

const final = result.cut().encode();

// Step 4: Send the Data to the Printer
endpoint.transfer(final, (error) => {
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
