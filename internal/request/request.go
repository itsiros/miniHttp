package request

import (
	"fmt"
	"io"
	// "os"
	"strings"
	"unicode"
)

const bufferSize int = 8

type parseState int

const (
	INITIALIZED = iota
	DONE
)

type Request struct {
	RequestLine RequestLine
	parseState  parseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func isKeywordCapitalized(key string) bool {
	for _, c := range key {
		if unicode.IsLetter(c) && !unicode.IsUpper(c) {
			return false
		}
	}
	return true
}

func (r *Request) parse(data []byte) (int, error) {

	if r.parseState == DONE {
		return -1, fmt.Errorf("error: trying to read data in DONE state")
	}
	if r.parseState != INITIALIZED {
		return -2, fmt.Errorf("error: unknown state")
	}

	bytesRead, err := parseRequestLine(string(data), r)
	if err != nil {
		return -3, err
	}

	if bytesRead == 0 {
		return 0, nil
	}

	r.parseState = DONE
	return bytesRead, nil
}

func parseRequestLine(line string, r *Request) (int, error) {

	idx := strings.Index(line, "\r\n")
	if idx == -1 {
		return 0, nil
	}

	requestLine := line[:idx]

	parts := strings.Fields(requestLine)

	if len(parts) != 3 {
		return 0, fmt.Errorf("Invalid request: %q", requestLine)
	}

	method, target, version := parts[0], parts[1], parts[2]

	if !isKeywordCapitalized(method) || !isKeywordCapitalized(version) {
		return 0, fmt.Errorf("version or method are not capitalized: %s", line)
	}

	if method != "GET" && method != "POST" && method != "HEAD" {
		return 0, fmt.Errorf("wrong method or non existand: %s", method)
	}

	if !strings.HasPrefix(target, "/") {
		return 0, fmt.Errorf("Invalid target: %s", target)
	}

	if version != "HTTP/1.1" {
		return 0, fmt.Errorf("Wrong HTTP version requested: %s", version)
	}

	r.RequestLine.Method = method
	r.RequestLine.RequestTarget = target
	r.RequestLine.HttpVersion = "1.1"

	return idx + 2, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, bufferSize)
	readToIndex := 0

	r := &Request{parseState: INITIALIZED}

	for r.parseState != DONE {

		if readToIndex == len(buf) {
			newBuf := make([]byte, 2*len(buf))
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			return nil, err
		}

		if err == io.EOF {
			r.parseState = DONE
			break
		}

		readToIndex += n
		consumed, parseErr := r.parse(buf[:readToIndex])
		if parseErr != nil {
			return nil, parseErr
		}

		if consumed > 0 {
			copy(buf, buf[consumed:readToIndex])
			readToIndex -= consumed
		}
	}
	return r, nil
}

// 	request := strings.parts(string(buf), "\r\n")
// 	if len(request) == 0 || request[0] == "" {
// 		return nil, fmt.Errorf("empty request")
// 	}
//
// 	info, parseErr := parseRequestLine(request[0])
// 	if parseErr != nil {
// 		fmt.Fprintln(os.Stderr, parseErr)
// 		return nil, parseErr
// 	}
//
// 	baby := Request{
// 		RequestLine: RequestLine{
// 			Method:        info[0],
// 			RequestTarget: info[1],
// 			HttpVersion:   info[2][5:],
// 		},
// 	}
//
// 	return &baby, nil
// }
