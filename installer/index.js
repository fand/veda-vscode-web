#!/usr/bin/env node
const fs = require('fs');
const path = require('path');
const sudo = require('sudo-prompt');

const exePath = process.argv[0];
const exeDir = path.dirname(exePath);
const filename = 'gl.veda.vscode.web.server.json';
const srcPath = path.resolve(exeDir, `../Resources/${filename}`);
const dstPath = `/Library/Google/Chrome/NativeMessagingHosts/${filename}`;

const iconPath = path.join(__dirname, 'icon.icns');
const options = {
  name: 'VEDA for VSCode Web Server',
  icns: iconPath,
};


if (fs.existsSync(dstPath)) {
  process.exit(0);
} else {
  sudo.exec(`cp ${srcPath} ${dstPath}`, options, (error) => {
    process.exit(error ? -1 : 0);
  });
}

