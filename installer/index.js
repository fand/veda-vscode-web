#!/usr/bin/env node
const fs = require('fs');
const path = require('path');
const sudo = require('sudo-prompt');
const cp = require('child_process');

const exePath = process.argv[0];
const exeDir = path.dirname(exePath);

const iconPath = path.join(__dirname, 'icon.icns');
const options = {
  name: 'VEDA for VSCode Web Server',
  icns: iconPath,
};

const filename = 'gl.veda.vscode.web.server.json';
const srcPath = path.resolve(exeDir, `../Resources/${filename}`);
const dstPath = `~/Library/Application\ Support/Google/Chrome/NativeMessagingHosts/${filename}`;

try {
  const dp = dstPath.replace(/(\s+)/g, '\\$1');
  console.log(dp);
  cp.execSync(`cat ${dp}`) // Fails if not exist
} catch (e) {
  const sp = srcPath.replace(/(\s+)/g, '\\$1');
  const dp = dstPath.replace(/(\s+)/g, '\\$1');
  console.log(`cp ${sp} ${dp}`);
  sudo.exec(`cp ${sp} ${dp}`, options, (error) => {
    if (error) { process.exit(-1); }
  });
}

