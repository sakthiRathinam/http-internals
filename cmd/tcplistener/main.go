package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net"

	request "github.com/sakthiRathinam/http_internals/internal/request"
)

func main() {
	printReqBytes := flag.Bool("printbytes", false, "to print the req bytes")
	flag.Parse()
	fmt.Println(*printReqBytes)
	tcpListener, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		fmt.Println("Some error while creating the tcp server", err)
		return
	}
	defer tcpListener.Close()
	for {
		connection, err := tcpListener.Accept()
		if err != nil {
			fmt.Println("Err happened while accpeting the connection", err)
			continue
		}
		if *printReqBytes {
			fmt.Println("Printing the bytes")
			printBytes(connection)
			connection.Close()
			continue
		}
		parsedReqObj, err := request.RequestFromReader(connection)
		if err != nil {
			fmt.Println("error occurend while parsing", err)
			return
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", parsedReqObj.RequestLine.Method)
		fmt.Println("- Target:", parsedReqObj.RequestLine.RequestTarget)
		fmt.Println("- Version:", parsedReqObj.RequestLine.HttpVersion)
	}

}

func printBytes(conn io.Reader) {
	byteString := ""
	for {
		buffer := make([]byte, 120)
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				byteString += string(buffer[:n])
				break
			}
			fmt.Println("Error while reading the bytes", err)
			return
		}
		byteString += string(buffer[:n])
	}
	// 1) Show escapes (\r\n):
	fmt.Printf("Raw (%d bytes): %q\n", len(byteString), byteString)

	// 2) Or, show a hex + ASCII dump:
	fmt.Printf("Hex dump:\n%s\n", hex.Dump([]byte(byteString)))
}
