package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	var port int
	var addr string

	// configure ip, port by setting the executive arguments
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.StringVar(&addr, "address", "localhost", "-address localhost")
	flag.Parse()

	// tcp connection using net package
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal("ResolveTCPAddr failed: ", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal("Dial failed: ", err)
	}

	// print data through conn
	go func() {
		scan := bufio.NewScanner(conn)
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			fmt.Println(scan.Text())
		}
	}()

	// send text input from the keyboard
	for {
		inputScan := bufio.NewScanner(os.Stdin)
		inputScan.Split(bufio.ScanLines)
		for inputScan.Scan() {
			if inputScan.Text() == "exit" {
				return
			}
			conn.Write([]byte(fmt.Sprintf("%s\n", inputScan.Text())))
		}
	}
}
