package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func main() {

	//create a tcp listener on port 80 
	listener, err := net.Listen("tcp", ":8080")

	if err != nil { //err handling
		fmt.Fprintf(os.Stderr , "Error startin server: %v\n" , err)
		os.Exit(1) // exit the program
	}

	defer listener.Close(); //Stop serving when all is done

	fmt.Println("Smuggling server running on port 8080...")
	fmt.Println("Waiting for clients to request files...")

	for { 

		conn , err := listener.Accept()  //block and connect
			if err != nil { 
				fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
				continue 
			}

		//connection succeded now Clint <======>Server
		fmt.Println("New client connected")
		go handleClient(conn) //hanlde multpiple clients at the same time in sep go routines.
	}
}

func handleClient(conn net.Conn) { 

	defer conn.Close()//close the connection at the end 

	var nameLen uint32

	err := binary.Read(conn , binary.LittleEndian, &nameLen)

	if err != nil { 
		fmt.Printf("Error reading file name len: %v\n", err)
		return
	}

	nameBuf := make([]byte , nameLen)

	_ ,err = conn.Read(nameBuf)
	if err != nil { 
		fmt.Printf("Error reading file name: %v\n", err)
		return
	}

	fmt.Println("Client requested")

	fileData , err := os.ReadFile(string(nameBuf))

	if err != nil {
		fmt.Println("Error reading file: ", err)

		binary.Write(conn, binary.LittleEndian, uint32(0)) //idk what this does yet
		return
	}

	fmt.Printf("Sending %d bytes...\n", len(fileData))


	// Send filename length + filename
	binary.Write(conn, binary.LittleEndian, uint32(len(fileData)))
	conn.Write([]byte(fileData))

	// Send file data length + data
	binary.Write(conn, binary.LittleEndian, uint32(len(fileData)))
	conn.Write(fileData)


	fmt.Println("File sent succesfully!")
}
