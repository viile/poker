package main

import (
	"flag"
	"fmt"
	"github.com/viile/poker/server"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var (
		host string
		port string
	)
	flag.StringVar(&host, "h", "0.0.0.0", "server host")
	flag.StringVar(&port, "p", "8787", "server port")
	flag.Parse()

	addr := fmt.Sprintf("%s:%s", host, port)

	s, err := server.NewServer(addr)
	if err != nil {
		panic(err)
	}

	s.Run()
}
