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

func parseRequestLine(line string, baby Request) error {

	if line == "" {
		return fmt.Errorf("Empty string")
	}

	split := strings.Split(line, " ")
	fmt.Printf("ITS ME: '%s' , '%s', '%s'\n", split[0], split[1], split[2])

	if !isKeywordCapitalized(split[0]) || !isKeywordCapitalized(split[2]) {
		return fmt.Errorf("Not all letters are capitalized: %s", line)
	}

	if strings.Contains(split[0], "GET") || strings.Contains(split[0], "POST") {
		baby.RequestLine.Method = split[0]
	}

	if split[1][0] == '/' {
		baby.RequestLine.RequestTarget = split[1]
	}

	if strings.Contains(split[2], "HTTP/") {
		baby.RequestLine.HttpVersion = split[2][5:]
	}

	return nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	var baby Request

	buf, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	request := strings.Split(string(buf), "\r\n")
	parseErr := parseRequestLine(request[0], baby)
	if parseErr != nil {
		fmt.Fprintln(os.Stderr, "Error: ", parseErr)
	}

	fmt.Print(baby.RequestLine.HttpVersion, baby.RequestLine.RequestTarget, baby.RequestLine.Method)
	return &baby, err
}
