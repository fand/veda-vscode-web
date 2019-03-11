package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unsafe"

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
		panic("Could not determine native endianness.")
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
		panic(err)
	}

	// Parse size
	size := nativeEndian.Uint32(sizeBuf)

	// Read content
	bufContent := make([]byte, size)
	_, err = io.ReadFull(os.Stdin, bufContent)
	if err != nil {
		panic(err)
	}

	var request request
	if err := json.Unmarshal(bufContent, &request); err != nil {
		panic(err)
	}

	return request
}

func sendResponse(res response) {
	str, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	println(str)
}

func launchServer() int {
	port := -1

	// Return the port if the server is already running
	processes, err := exec.Command("ps aux").Output()
	if err != nil {
		log.Fatal(err)
	}
	for _, process := range strings.Split(string(processes), "\n") {
		if strings.Index(process, "code-server-wrapper") != -1 {
			args := strings.Split(process, " ")
			portStr := args[len(args)-1]
			port, err = strconv.Atoi(portStr)
			if err != nil {
				panic(err)
			}
		}
	}

	// Launch code-server-wrapper if not running
	if port == -1 {
		port, err = freeport.GetFreePort()
		if err != nil {
			log.Fatal(err)
		}

		cmd := exec.Command("code-server-wrapper " + string(port))
		err = cmd.Start()
		if err != nil {
			panic(err)
		}
	}

	return port
}

func getServerPort(wrapperPort int) int {
	res, err := http.Get("http://localhost:" + strconv.Itoa(wrapperPort) + "/port")
	if err != nil {
		log.Fatal(err)
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	port, err := strconv.Atoi(string(resData))
	if err != nil {
		log.Fatal(err)
	}

	return port
}

func main() {
	wrapperPort := launchServer()
	serverPort := getServerPort(wrapperPort)

	// Open browser
	if len(os.Args) == 0 {
		open.Start("http://localhost:" + strconv.Itoa(serverPort))
	}

	// Handle messages from the extension via Native Messsaging API
	for {
		req := getRequest()

		fileURI := strings.Replace(req.fileUri, "file:///", "/", 1)

		bytes, err := ioutil.ReadFile(fileURI)
		if err != nil {
			panic(err)
		}

		res := response{shader: string(bytes)}
		sendResponse(res)
	}

}
