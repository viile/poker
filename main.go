package main

import (
	"flag"
	"fmt"
	"github.com/viile/poker/server"
	"log"
)

func main() {
	var (
		host string
		port string
		game string
		err error
	)
	flag.StringVar(&host, "h", "0.0.0.0", "server host")
	flag.StringVar(&port, "p", "8787", "server port")
	flag.StringVar(&game, "g", "poker", "server game")
	flag.Parse()

	addr := fmt.Sprintf("%s:%s",host,port)

	var s *server.Server
	s,err = server.NewServer(addr)
	if err != nil {
		log.Panic(err)
		return
	}
	s.Run()
}
