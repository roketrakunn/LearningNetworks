package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)


func main() { 
	conn , err := net.Dial("tcp", "Localhost:8080")

	if err != nil { 
		fmt.Fprintf(os.Stderr, "Error connecting :v%\n",err)
		fmt.Fprintf(os.Stderr, "Make sure the server is running!\n")
		os.Exit(1)
	}

	defer conn.Close() // close the connection when done 

	fmt.Println("Connected to echo server")
	

	//send the message
	message := "Hello echo sever.\n"
	fmt.Printf("Sending message :%s" , message)

	_ , err = conn.Write([]byte(message)) 

	if err != nil { 
		fmt.Fprintf(os.Stderr , "Error writing: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(conn)

	response , err := reader.ReadString('\n')
	
	if err != nil { 
		fmt.Fprintf(os.Stderr , " Error reading: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Recieved: %s" , response)
	fmt.Println("DONE!")
}
