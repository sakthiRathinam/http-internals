package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
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
		fmt.Println("connection has established")
		linesChan := getLinesChannel(connection)
		for line := range linesChan {
			fmt.Print(line)
		}
		fmt.Println("Connection has been closed")
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	var currentLine string
	linesChan := make(chan string, 10)
	go func() {
		defer f.Close()
		defer close(linesChan)
		for {
			byteBuffer := make([]byte, 8)
			n, err := f.Read(byteBuffer)
			if err != nil {
				if currentLine != "" {
					linesChan <- fmt.Sprintf("read: %s\n", currentLine)
					currentLine = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Println("Error: ", err)
				break
			}
			currentBuffStr := string(byteBuffer[:n])

			parts := strings.Split(currentBuffStr, "\n")
			for i := 0; i < len(parts)-1; i++ {
				linesChan <- fmt.Sprintf("read: %s%s\n", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()
	return linesChan
}
