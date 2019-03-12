#!/usr/bin/env node
const fs = require('fs');
const path = require('path');
const sudo = require('sudo-prompt');

const exePath = process.argv[0];
const exeDir = path.dirname(exePath);

const iconPath = path.join(__dirname, 'icon.icns');
const options = {
  name: 'VEDA for VSCode Web Server',
  icns: iconPath,
};

const filename = 'gl.veda.vscode.web.server.json';
const srcPath = path.resolve(exeDir, `../Resources/${filename}`).replace(/(\s+)/g, '\\$1');
const dstPath = `~/Library/Application\ Support/Google/Chrome/NativeMessagingHosts/${filename}`.replace(/(\s+)/g, '\\$1');

if (fs.existsSync(dstPath)) {
  process.exit(0);
} else {
  sudo.exec(`cp ${srcPath} ${dstPath}`, options, (error) => {
    if (error) { process.exit(-1); }
  });
}

