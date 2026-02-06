/**
 	+ This is the sever part that will echo whatever the client sends
	+ Basically like that one annoying friend that repeats whatever you say to them 
	+ What does this do ? 
			++ accepts a connection 
			++ reads  whatever the client sent 
			++ Yells it back 
			++ for only this case ... i just serve one client at a time 
	+ What this teaches me ? 
			++ How simple servers are I guess
*/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	/**
		+ Create a new tcp sever on port 80(http)
		+ This is what stays alive while there is more clients to serve
		+ Unlike connfd , This doesn not close whilst the program is still running
		+ net.Listen() returns the listenfd an an error if anything bad happened 
		+ We in this case need just the file descriptor more than the error i mean who likes errors
	*/
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening: %v\n", err)
		os.Exit(1)
	}
	/**
		+ This is what i was talking about
		+ This HAS to stay alive untill the program ends 
		+ Its like a resturant that is open all day serving clients till 6PM
		+ This is by far the best analogy for this stuff 
	*/
	defer listener.Close()
	
	fmt.Println("Echo server listening on port 8080...")
	fmt.Println("Waiting for clients to connect...")
	
	/**
		+ There it is .. see ? THis stays open and accepting clients 
		+ This is the resturant dude. 
	*/	
	for {
		/**
			+ net listener.Accept() returns a file descriptor and an error if anything goes wrong.
			+ This is the file descriptor we will do the read and writes to 
			+ Below is like a waiter serving you uhh pizza i guess .. 
			+ Notice how you will leave and another customer comes in ? 
			+ Yes that is the beauty of it man.
		*/
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting: %v\n", err)
			continue // Skip this client, keep listening
		}
		
		fmt.Printf("Client connected from: %s\n", conn.RemoteAddr())
		
		//handle whatever you want .. in this case  ist just reading and wriing back whwat you said
		handleClient(conn)
	}
}

/**
	+ This handles clientX ...
	+ Creates a read buffer 
	+ Reads the messsage from the client 
	+ Save this to response 
	+ add "ECHO" to the response
	+ writes it back to the handleClient
	+ And that my boy is why people yell io connfd .. if you do not get it then...idk 
*/
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
