package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unsafe"
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

	// Parse content
	// content := string(bufContent)

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

func main() {
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
