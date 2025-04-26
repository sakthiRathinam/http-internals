package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"
const colon = ":"

func validateHeader(headerVals []string) error {
	if len(headerVals) <= 1 {
		return fmt.Errorf("Not a valid header")
	}
	headerKey := headerVals[0]
	validKey := strings.Contains(headerKey, " ")
	if validKey {
		return fmt.Errorf("Not a valid header")
	}
	return nil
}
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	splitFieldLines := bytes.Split(data, []byte(crlf))
	if len(splitFieldLines) == 0 {
		return 0, false, err
	}

	fieldLine := string(splitFieldLines[0])
	if len(fieldLine) == 0 {
		return 0, true, nil
	}
	parsedFieldLineObj := strings.Split(fieldLine, colon)
	err = validateHeader(parsedFieldLineObj)
	if err != nil {
		return n, false, err
	}
	h[parsedFieldLineObj[0]] = strings.ReplaceAll(strings.Join(parsedFieldLineObj[1:], ":"), " ", "")
	return n, done, err
}
