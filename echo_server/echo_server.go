package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// 1. Listen on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close() 
	
	fmt.Println("Echo server listening on port 8080...")
	fmt.Println("Waiting for clients to connect...")
	
	// Server loop: keep accepting clients forever
	for {
		// 2. Accept a connection (blocks until client connects)
		// Returns a connection object (like a fd)
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting: %v\n", err)
			continue // Skip this client, keep listening
		}
		
		fmt.Printf("Client connected from: %s\n", conn.RemoteAddr())
		
		// 3. Handle this client
		handleClient(conn)
	}
}

// handleClient reads from client and echoes back
func handleClient(conn net.Conn) {
	defer conn.Close() //close connection when done
	
	// Create a buffered reader
	reader := bufio.NewReader(conn)
	
	// Read one line from client
	// ReadString reads until it sees '\n'
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading: %v\n", err)
		return
	}
	
	fmt.Printf("Received: %s", message) 
	
	// Echo it back with "ECHO: " prefix
	response := "ECHO: " + message
	
	// Write response back to client
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing: %v\n", err)
		return
	}
	
	fmt.Println("Response sent. Closing connection.")
}
