package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	. "github.com/fand/veda-vscode-web/server/logger"
	"github.com/phayes/freeport"
	"github.com/skratchdot/open-golang/open"
)

func getEndian() binary.ByteOrder {
	var nativeEndian binary.ByteOrder

	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		nativeEndian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		nativeEndian = binary.BigEndian
	default:
		LogFatal("Could not determine native endianness.")
	}

	return nativeEndian
}

type request struct {
	fileUri string
}

type response struct {
	shader string
}

func getRequest() request {
	var nativeEndian = getEndian()

	// Read size
	sizeBuf := make([]byte, 4)
	_, err := io.ReadFull(os.Stdin, sizeBuf)
	if err != nil {
		LogFatal(err)
	}

	// Parse size
	size := nativeEndian.Uint32(sizeBuf)

	// Read content
	bufContent := make([]byte, size)
	_, err = io.ReadFull(os.Stdin, bufContent)
	if err != nil {
		LogFatal(err)
	}

	var request request
	if err := json.Unmarshal(bufContent, &request); err != nil {
		LogFatal(err)
	}

	return request
}

func sendResponse(res response) {
	str, err := json.Marshal(res)
	if err != nil {
		LogFatal(err)
	}
	println(str)
}

func launchServerWrapper() int {
	port := -1

	// Get file paths
	exePath, err := os.Executable()
	if err != nil {
		LogFatal(err)
	}
	dirPath := filepath.Dir(exePath)
	wrapperCmdPath := filepath.Join(dirPath, "code-server-wrapper")

	// Return the port if the server is already running
	processes, err := exec.Command("ps", "aux").Output()
	if err != nil {
		LogFatal(err)
	}
	for _, process := range strings.Split(string(processes), "\n") {
		if strings.Index(process, "code-server-wrapper") != -1 {
			args := strings.Split(process, " ")
			portStr := args[len(args)-1]
			port, err = strconv.Atoi(portStr)
			if err != nil {
				LogFatal(err)
			}
		}
	}

	// Launch code-server-wrapper if not running
	if port == -1 {
		port, err = freeport.GetFreePort()
		if err != nil {
			LogFatal(err)
		}

		cmd := exec.Command(wrapperCmdPath, strconv.Itoa(port))
		err = cmd.Start()
		if err != nil {
			LogFatal(err)
		}
	}

	return port
}

func httpGet(url string, retry int) string {
	client := &http.Client{Timeout: 10 * time.Second}

	var res *http.Response
	var err error

	for i := 0; i < retry; i++ {
		res, err = client.Get(url)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		LogFatal(err)
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		LogFatal(err)
	}

	return string(resData)
}

func getServerPort(wrapperPort int) int {
	res := httpGet(fmt.Sprintf("http://127.0.0.1:%d/port", wrapperPort), 10)

	port, err := strconv.Atoi(res)
	if err != nil {
		LogFatal(err)
	}

	return port
}

func isOpenFromFinder() bool {
	for _, arg := range os.Args {
		if strings.Index(arg, "-psn") != -1 {
			return true
		}
	}
	return false
}

func cleanup() {
	FlushLogger()
}

func main() {
	InitLogger("/tmp/gl.veda.vscode.web.server/log-main.txt")

	// Prepare cleanup
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-chSignal
		cleanup()
	}()

	// Launch code-server-wrapper
	wrapperPort := launchServerWrapper()
	serverPort := getServerPort(wrapperPort)

	// Open browser when the app is launched from Finder
	if isOpenFromFinder() {
		time.AfterFunc(3*time.Second, func() {
			open.Start(fmt.Sprintf("http://localhost:%d", serverPort))
		})
	}

	// Handle messages from the extension via Native Messsaging API
	for {
		req := getRequest()

		fileURI := strings.Replace(req.fileUri, "file:///", "/", 1)

		bytes, err := ioutil.ReadFile(fileURI)
		if err != nil {
			LogFatal(err)
		}

		res := response{shader: string(bytes)}
		sendResponse(res)
	}

	cleanup()
}
