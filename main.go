package main
import (
	"flag"
	"fmt"
	"log"

	gnet "github.com/panjf2000/gnet/v2"
)

type echoServer struct {
	gnet.BuiltinEventEngine

	eng gnet.BuiltinEventEngine
	addr string
	multicore bool
}