package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	d, err := net.Dial("tcp4", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		if _, err = d.Write([]byte(text)); err != nil {
			log.Fatalln(err)
		}
		text, err = bufio.NewReader(d).ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(text))
	}
}
