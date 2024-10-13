import EscPosEncoder from 'esc-pos-encoder';
import { findByIds } from 'usb';

import { createCanvas, loadImage } from 'canvas'

const canvas = createCanvas(200, 200)
const ctx = canvas.getContext('2d')

// Draw cat with lime helmet
loadImage('./image.jpg').then((image) => {
  ctx.drawImage(image, 50, 0, 70, 70)
})

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
  .image(canvas, 400, 400).encode()

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
