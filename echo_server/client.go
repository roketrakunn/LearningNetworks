/**
 + This is a basic client that sends a message to a local server
 + It then reads the response from the client and prints to the screen 
 + What it teaches me ? How simple web stuff is... definetley overrated
 + The server sends a message back 
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
	+ tcp connection  localhost on port 80(http)
	+ if err is not null , the connection has failed 
	+ Confirm sever is up.
	*/
	conn , err := net.Dial("tcp", "Localhost:8080")

	if err != nil { 
		fmt.Fprintf(os.Stderr, "Error connecting :v%\n",err)
		fmt.Fprintf(os.Stderr, "Make sure the server is running!\n")
		os.Exit(1)
	}

	//at end of main close the connection(just good practice)
	defer conn.Close() 

	fmt.Println("Connected to echo server")
	

	//send the message
	message := "Hello echo sever.\n" //message you uwanna send
	fmt.Printf("Sending message :%s" , message)

	_ , err = conn.Write([]byte(message)) // write ro the file ( see I told you everhting is a file)
	if err != nil { 
		fmt.Fprintf(os.Stderr , "Error writing: %v\n", err)
		os.Exit(1)
	}

	/**
	+ Create a buffer reader
	+ Will read the response from server
	+ In this case the resposne will be "ECHO: <message>"
	*/
	reader := bufio.NewReader(conn) 

	response , err := reader.ReadString('\n')  // save response to response
	
	if err != nil { 
		fmt.Fprintf(os.Stderr , " Error reading: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Recieved: %s" , response) //print the response to your terminal ...
	fmt.Println("DONE!")
}
