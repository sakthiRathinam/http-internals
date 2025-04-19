package main

import (
	"fmt"
	"net"

	request "github.com/sakthiRathinam/http_internals/internal/request"
)

func main() {
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
