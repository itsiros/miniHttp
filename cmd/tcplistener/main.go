package main

import (
	"fmt"
	"github.com/tsironi93/miniHttp/internal/request"
	"net"
)

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

		req, err := request.RequestFromReader(fd)
		if err != nil {
			fmt.Println("parse error:", err)
		} else {
			fmt.Printf("parse request: %+v\n", req.RequestLine)
		}

		fd.Close()
	}
}
