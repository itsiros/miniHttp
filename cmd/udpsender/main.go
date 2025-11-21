package main

import (
	// "fmt"
	"net"
)

func main() {

	udp, er := net.ResolveUDPAddr("udp", "localhost:42069")
	if er == nil {
		panic(er)
	}

	net.DialUDP("udp")
}
