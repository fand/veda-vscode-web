let port;

chrome.browserAction.onClicked.addListener(tab => {
  if (!port) {
    port = chrome.runtime.connectNative("gl.veda.vscode.web.server");

    // Transfer request from content to server
    chrome.runtime.onMessage.addListener((req, sender, sendResponse) => {
      if (!req.fileUri) { return; }
      if (!port) { return; }
      port.postMessage(req);
    });

    // Transfer response from server to content
    port.onMessage.addListener((res) => {
      console.log('>> received', res);
      chrome.tabs.sendMessage(tab.id, res);
    });
  }

  chrome.tabs.executeScript(tab.id, {
    file: 'content.js'
  }, () => {
  });
});

