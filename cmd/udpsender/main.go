package main

import (
	// "fmt"
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	udp, er := net.ResolveUDPAddr("udp", "localhost:42069")
	if er != nil {
		panic(er)
	}

	udpfd, err := net.DialUDP("udp", nil, udp)
	if err != nil {
		panic(err)
	}
	defer udpfd.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if input == "" {
			continue
		}

		_, er := udpfd.Write([]byte(input))
		if er != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
}
