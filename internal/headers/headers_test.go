package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadersParse(t *testing.T) {
	testCases := []struct {
		headerByte []byte
		headers    []struct {
			headerKey string
			headerVal string
			errString string
		}
	}{
		{
			headerByte: []byte("Host: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"),
			headers: []struct {
				headerKey string
				headerVal string
				errString string
			}{
				{
					headerKey: "Host",
					headerVal: "localhost:42069",
					errString: "",
				},
				{
					headerKey: "User-Agent",
					headerVal: "curl/7.81.0",
					errString: "",
				},
				{
					headerKey: "Accept",
					headerVal: "*/*",
					errString: "",
				},
			},
		}, // valid header test case
		{
			headerByte: []byte("Host : localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"),
			headers: []struct {
				headerKey string
				headerVal string
				errString string
			}{
				{
					headerKey: "Host",
					headerVal: "localhost:42069",
					errString: "not a valid header",
				},
			},
		}, // Invalid spacing header key
		{
			headerByte: []byte("Host: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"),
			headers: []struct {
				headerKey string
				headerVal string
				errString string
			}{
				{
					headerKey: "Host",
					headerVal: "localhost:42069",
					errString: "",
				},
				{
					headerKey: "User-Agent",
					headerVal: "curl/7.81.0",
					errString: "",
				},
				{
					headerKey: "Accept",
					headerVal: "*/*",
					errString: "",
				},
			},
		},
	}

	headers := Headers{}
	done := false
	n := 0
	for _, testCase := range testCases {
		headerIdx := 0
		for done == false {
			parsedBytes, endLoop, err := headers.Parse(testCase.headerByte[n:])
			if endLoop {
				break
			}
			testCaseHeader := testCase.headers[headerIdx]
			if testCaseHeader.errString != "" {
				assert.Equal(t, testCaseHeader.errString, err.Error())
			}
			if testCaseHeader.headerKey != "" {
				assert.Equal(t, testCaseHeader.headerVal, headers[testCaseHeader.headerKey])
			}

			n += parsedBytes
			headerIdx++
		}
	}
}
