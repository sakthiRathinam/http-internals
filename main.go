package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	messageTxtFileObj, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error opening message.txt:", err)
		return
	}
	defer messageTxtFileObj.Close()
	linesChan := getLinesChannel(messageTxtFileObj)
	for line := range linesChan {
		fmt.Print(line)
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
