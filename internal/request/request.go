package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method        string
	HttpVersion   string
	RequestTarget string
}

func RequestFromReader(reader io.Reader) (req *Request, err error) {
	wholeReqBytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("some error happend while reading the request obj", err)
		return
	}
	splitReqByLines := strings.Split(string(wholeReqBytes), "\r\n")
	requestLineObjs := strings.Split(splitReqByLines[0], " ")
	fmt.Println(requestLineObjs)
	if len(requestLineObjs) != 3 {
		return &Request{}, errors.New("invalid request")
	}
	if requestLineObjs[2] != "HTTP/1.1" {
		return &Request{}, errors.New("we only support HTTP/1.1")
	}
	requestLineObj := RequestLine{requestLineObjs[0], "1.1", requestLineObjs[1]}
	return &Request{requestLineObj}, nil
}
