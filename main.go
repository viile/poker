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
		mode string
		err error
	)
	flag.StringVar(&host, "h", "0.0.0.0", "server host")
	flag.StringVar(&port, "p", "8787", "server port")
	flag.StringVar(&mode, "m", "server", "server mode")
	flag.Parse()

	addr := fmt.Sprintf("%s:%s",host,port)

	switch mode {
	case "server":
		var s *server.Server
		s,err = server.NewServer(addr)
		if err != nil {
			log.Panic(err)
			return
		}
		s.Run()
	default:

	}

}
