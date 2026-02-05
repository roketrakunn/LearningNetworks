package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

// FileData represents the data we receive from the server
type FileData struct {
	Name string
	Data []byte
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client.go <filename>")
		fmt.Println("Example: go run client.go secret_data.txt")
		return
	}

	requestedFile := os.Args[1]

	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to smuggling server!\n")
	fmt.Printf("Requesting: %s\n", requestedFile)

	// Send the filename we want (length + name)
	binary.Write(conn, binary.LittleEndian, uint32(len(requestedFile)))
	conn.Write([]byte(requestedFile))

	// Receive the response
	fileData, err := receiveFile(conn)
	if err != nil {
		fmt.Println("Error receiving file:", err)
		return
	}

	fmt.Printf("Received: %s (%d bytes)\n", fileData.Name, len(fileData.Data))
	fmt.Println("\n--- FILE CONTENTS ---")
	fmt.Println(string(fileData.Data))
	fmt.Println("--- END ---")

	
	outputName := "downloaded_" + fileData.Name
	os.WriteFile(outputName, fileData.Data, 0644)
	fmt.Printf("\nSaved to: %s\n", outputName)
}

func receiveFile(conn net.Conn) (*FileData, error) {
	// Read filename length
	var nameLen uint32
	err := binary.Read(conn, binary.LittleEndian, &nameLen)
	if err != nil {
		return nil, err
	}

	if nameLen == 0 {
		return nil, fmt.Errorf("server couldn't find the file")
	}

	// Read filename
	nameBuf := make([]byte, nameLen)
	_, err = conn.Read(nameBuf)
	if err != nil {
		return nil, err
	}

	// Read file data length
	var dataLen uint32
	err = binary.Read(conn, binary.LittleEndian, &dataLen)
	if err != nil {
		return nil, err
	}

	// Read file data
	dataBuf := make([]byte, dataLen)
	_, err = conn.Read(dataBuf)
	if err != nil {
		return nil, err
	}

	return &FileData{
		Name: string(nameBuf),
		Data: dataBuf,
	}, nil
}
