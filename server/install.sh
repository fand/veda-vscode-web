#!/bin/sh
set -eux

DIR="$( cd "$( dirname "$0" )" && pwd )"
if [ "$(uname -s)" = 'Darwin' ]; then
  TARGET_DIR='/Library/Google/Chrome/NativeMessagingHosts'
else
  TARGET_DIR='/etc/opt/chrome/native-messaging-hosts'
fi

HOST_NAME=gl.veda.vscode.web.server

# Create directory to store native messaging host.
mkdir -p $TARGET_DIR

# Copy native messaging host manifest.
cp $DIR/$HOST_NAME.json $TARGET_DIR

# Update host path in the manifest.
HOST_PATH=$DIR/server.py
ESCAPED_HOST_PATH="$(printf '%s\n' "$HOST_PATH" | sed 's/#/\\#/g')"
sed -i -e "s#HOST_PATH#$ESCAPED_HOST_PATH#" $TARGET_DIR/$HOST_NAME.json

# Set permissions for the manifest so that all users can read it.
chmod o+r $TARGET_DIR/$HOST_NAME.json

echo "VEDA for VSCode Web Server has been installed."
