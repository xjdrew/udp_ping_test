package main

import (
	"flag"
	"log"
	"net"
)

func main() {
	listen := flag.String("listen", ":8888", "listen addr")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *listen)
	if err != nil {
		log.Fatal("net.ResolveUDPAddr():", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("net.ListenUDP():", err)
	}
	defer conn.Close()

	log.Printf("listen on %s", *listen)

	b := make([]byte, 2048)
	for {
		cc, remote, rderr := conn.ReadFromUDP(b)

		if rderr != nil {
			log.Printf("read packet failed: %s", rderr)
			break
		}

		i := 8
		if cc < i {
			i = cc
		}

		log.Printf("[%v]<%s> read %d bytes", remote, string(b[:i]), cc)
		cc, wrerr := conn.WriteTo(b[:cc], remote)
		if wrerr != nil {
			log.Printf("write packet failed: %s\n", wrerr)
			break
		}
	}

	log.Print("done")
}
