package request

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data           string
	pos            int
	numBytePerRead int
}

func testBufferAppend(readObj *chunkReader) {
	buf := make([]byte, 100)
	readToIndex := 3
	_, err := readObj.Read(buf[readToIndex:])
	if err != nil {
		fmt.Println(err)
		return
	}
	readToIndex = 50
	fmt.Println(buf)
	_, err = readObj.Read(buf[readToIndex:])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(buf)
	readToIndex = 80
	_, err = readObj.Read(buf[readToIndex:])
	if err != nil {
		fmt.Println(err)
		return
	}
	readToIndex = 0
	fmt.Println(buf)
	_, err = readObj.Read(buf[readToIndex:])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(buf)
}
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytePerRead

	if endIndex > len(cr.data) {
		endIndex = len(cr.data) - 1
	}

	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytePerRead {
		cr.pos -= n - cr.numBytePerRead
		n = cr.numBytePerRead
	}
	return n, nil
}

func TestBufferBehaviour(t *testing.T) {
	reader := &chunkReader{
		data:           "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		pos:            0,
		numBytePerRead: 3,
	}
	testBufferAppend(reader)
}
func TestRequestLineParse(t *testing.T) {
	// Test: Good GET Request line
	reader := &chunkReader{
		data:           "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		pos:            0,
		numBytePerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	reader = &chunkReader{
		data:           "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		pos:            0,
		numBytePerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	reader = &chunkReader{
		data:           "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		pos:            0,
		numBytePerRead: 3,
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid request")

	// Test: Invalid HTTP version
	reader = &chunkReader{
		data:           "POST /coffee HTTP/1.2\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		pos:            0,
		numBytePerRead: 3,
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)
	require.ErrorContains(t, err, "we only support HTTP/1.1")
}
