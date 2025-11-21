package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {

	ch := make(chan string)

	go func() {
		defer close(ch)

		var line string
		for {

			buf := make([]byte, 8)
			n, er := io.ReadFull(f, buf)
			if er == io.EOF {
				break
			}

			line += string(buf[:n])
			if er == io.ErrUnexpectedEOF {
				ch <- line
				break
			}

			if er != nil {
				panic(er)
			}

			for {
				i := strings.IndexRune(line, '\n')
				if i == -1 {
					break
				}

				ch <- line[:i+1]
				line = line[i+1:]
			}
		}
	}()
	return ch
}

func main() {

	listener, er := net.Listen("tcp", "127.0.0.1:42069")
	if er != nil {
		panic(er)
	}
	defer listener.Close()

	for {
		fd, er := listener.Accept()
		if er != nil {
			panic(er)
		}

		fmt.Println("New connection established!")
		ch := getLinesChannel(fd)
		for chunk := range ch {
			fmt.Print(chunk)
		}
		fmt.Println("connection closed")
		fd.Close()
	}
}
