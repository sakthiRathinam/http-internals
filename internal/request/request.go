package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const crlf = "\r\n"
const bufferSize = 8

type Request struct {
	RequestLine RequestLine
	state       requestState
}
type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

type RequestLine struct {
	Method        string
	HttpVersion   string
	RequestTarget string
}

func RequestFromReader(reader io.Reader) (req *Request, err error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	reqObj := &Request{state: requestStateInitialized}
	for reqObj.state != requestStateDone {
		if readToIndex >= 0 {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				reqObj.state = requestStateDone
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead
		numBytesParsed, err := reqObj.Parse(buf[:readToIndex])
		if err != nil {
			if errors.Is(err, io.EOF) {
				reqObj.state = requestStateDone
				break
			}
			return reqObj, err
		}
		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return req, nil
}

func parseRequestLine(data []byte) (reqLine *RequestLine, n int, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLine, err := RequestLineFromStr(string(data[:idx]))
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}
func RequestLineFromStr(str string) (reqLine *RequestLine, err error) {
	splitReqByLines := strings.Split(str, crlf)
	requestLineObjs := strings.Split(splitReqByLines[0], " ")
	fmt.Println(requestLineObjs)
	if len(requestLineObjs) != 3 {
		return &RequestLine{}, errors.New("invalid request")
	}
	if requestLineObjs[2] != "HTTP/1.1" {
		return &RequestLine{}, errors.New("we only support HTTP/1.1")
	}
	requestLineObj := RequestLine{requestLineObjs[0], "1.1", requestLineObjs[1]}
	return &requestLineObj, nil
}

func (r *Request) Parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			// need more data for the request line
			return n, nil
		}
		r.RequestLine = *requestLine
		r.state = requestStateDone
		return n, nil
	case requestStateDone:
		return 0, fmt.Errorf("this request state was alredy done")
	default:
		return 0, fmt.Errorf("Unknown state error")
	}
}
