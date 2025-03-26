package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	gnet "github.com/panjf2000/gnet/v2"
)

// chat server sends data to all the connected clients
type chatServer struct {
	gnet.BuiltinEventEngine

	// Map that stores connections
	cliMap sync.Map
}

// OnOpen is calleed when a new client joins
func (cs *chatServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("client connected. address:%s", c.RemoteAddr().String())

	// stores new connection in Map
	cs.cliMap.Store(c, true)
	return nil, gnet.None
}

// OnClose is calleed when a client disconnects
func (cs *chatServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	log.Printf("client disconnected. address:%s", c.RemoteAddr().String())

	// deletes the dead connection from Map
	if _, ok := cs.cliMap.LoadAndDelete(c); ok {
		log.Printf("connection removed")
	}
	return gnet.None
}

// OnBoot is called when the server starts
func (cs *chatServer) OnBoot(eng gnet.Engine) gnet.Action {
	log.Printf("chat server is listening\n")
	return gnet.None
}

// OnTraffic is called when the server receives data from the network
func (cs *chatServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)

	// sends data to all connections upon receiving data
	cs.cliMap.Range(func(key, value any) bool {
		if conn, ok := key.(gnet.Conn); ok {
			conn.AsyncWrite(buf, nil)
		}
		return true
	})

	return gnet.None
}

func main() {
	var port int
	var multicore bool

	// Example command: go run echo.go --port 9000 --multicore=true
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore true")
	flag.Parse()

	chat := &chatServer{}

	log.Fatal(gnet.Run(chat, fmt.Sprintf("tcp://:%d", port), gnet.WithMulticore(multicore)))
}

// // echo server sends data to one client
// type echoServer struct {
// 	gnet.BuiltinEventEngine

// 	eng       gnet.Engine
// 	addr      string
// 	multicore bool
// }

// // OnOpen is calleed when a new client joins
// func (es *echoServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
// 	log.Printf("client connected. address:%s", c.RemoteAddr().String())
// 	return nil, gnet.None
// }

// // OnClose is calleed when a client disconnects
// func (es *echoServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
// 	log.Printf("client disconnected. address:%s", c.RemoteAddr().String())
// 	return gnet.None
// }

// // OnBoot is called when the server starts
// func (es *echoServer) OnBoot(eng gnet.Engine) gnet.Action {
// 	es.eng = eng
// 	log.Printf("echo server with multi-core=%t is listening on %s\n", es.multicore, es.addr)
// 	return gnet.None
// }

// // OnTraffic is called when the server receives data from the network
// func (es *echoServer) OnTraffic(c gnet.Conn) gnet.Action {
// 	buf, _ := c.Next(-1)
// 	c.Write(buf)
// 	return gnet.None
// }

// func main() {
// 	var port int
// 	var multicore bool

// 	// Example command: go run echo.go --port 9000 --multicore=true
// 	flag.IntVar(&port, "port", 9000, "--port 9000")
// 	flag.BoolVar(&multicore, "multicore", false, "--multicore true")
// 	flag.Parse()

// 	echo := &echoServer{
// 		addr:      fmt.Sprintf("tcp://:%d", port),
// 		multicore: multicore,
// 	}

// 	log.Fatal(gnet.Run(echo, echo.addr, gnet.WithMulticore(multicore)))
// }
