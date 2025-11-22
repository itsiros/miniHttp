package request

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
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

func parseRequestLine(line string) ([]string, error) {

	if line == "" {
		return nil, fmt.Errorf("Empty string")
	}

	split := strings.Fields(line)

	if len(split) != 3 {
		return nil, fmt.Errorf("Invalid request")
	}

	if !isKeywordCapitalized(split[0]) || !isKeywordCapitalized(split[2]) {
		return nil, fmt.Errorf("Not all letters are capitalized: %s", line)
	}

	if split[0] != "GET" && split[0] != "POST" && split[0] != "HEAD" {
		return nil, fmt.Errorf("No method specified")
	}

	if split[1][0] != '/' {
		return nil, fmt.Errorf("No request target")
	}

	if strings.Compare(split[2], "HTTP/1.1") != 0 {
		return nil, fmt.Errorf("Wrong HTTP version requested")
	}

	return split, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	request := strings.Split(string(buf), "\r\n")
	if len(request) == 0 || request[0] == "" {
		return nil, fmt.Errorf("empty request")
	}

	info, parseErr := parseRequestLine(request[0])
	if parseErr != nil {
		fmt.Fprintln(os.Stderr, parseErr)
		return nil, parseErr
	}

	baby := Request{
		RequestLine: RequestLine{
			Method:        info[0],
			RequestTarget: info[1],
			HttpVersion:   info[2][5:],
		},
	}

	return &baby, nil
}
