#!/usr/bin/env node
const sudo = require('sudo-prompt');
const path = require('path');

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

sudo.exec(`cp ${srcPath} ${dstPath}`, options, (error, stdout) => {
  console.error(error);
  console.log(stdout);
  console.log(error ? 'FAILURE' : 'SUCCESS');
});
