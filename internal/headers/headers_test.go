package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadersParse(t *testing.T) {
	testCases := []struct {
		headerByte []byte
		headerKey  string
		headerVal  string
		errString  string
	}{
		{[]byte("Host: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"),
			"Host", "localhost:42069",
			""}, // valid header test case
		{[]byte("Host : localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"),
			"Host", "localhost:42069",
			"Not a valid header"}, // Invalid spacing header key
		{[]byte("Host: localhost:42"),
			"", "", ""},

		{[]byte("Host: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"),
			"Host", "localhost:42069", "",
		},
	}

	headers := Headers{}
	for _, testCase := range testCases {
		_, _, err := headers.Parse(testCase.headerByte)
		if err != nil {

			assert.Equal(t, err.Error(), testCase.errString)
			continue
		}
		if testCase.headerKey != "" {
			assert.Equal(t, testCase.headerVal, headers[testCase.headerKey])
		}
	}
}
