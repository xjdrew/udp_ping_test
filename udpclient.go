package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	connect := flag.String("connect", "127.0.0.1:8888", "connect addr")
	sz := flag.Int("size", 512, "packet size")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *connect)
	if err != nil {
		log.Fatal("net.ResolveUDPAddr():", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("net.DialUDP:", err)
	}
	defer conn.Close()

	log.Printf("connect to %v", *connect)
	localAddr := conn.LocalAddr().String()

	if *sz <= 8 {
		*sz = 8
	}
	sbuf := make([]byte, *sz)
	rbuf := make([]byte, *sz)
	rand.Read(sbuf)

	index := 0
OuterLoop:
	for {
		salt := fmt.Sprintf("%08d", index)
		copy(sbuf, []byte(salt))
		index++

		if _, err := conn.Write(sbuf); err != nil {
			log.Printf("conn.Write():", err)
			break
		}

		log.Printf("[%s]<%s> write %d bytes", localAddr, salt, *sz)

		lost := true
		t0 := time.Now()
		for {
			// 1 second timeout
			conn.SetReadDeadline(t0.Add(time.Second))
			cc, rerr := conn.Read(rbuf)
			if rerr != nil {
				if opErr, ok := rerr.(*net.OpError); ok && opErr.Timeout() {
					// timeout
					break
				}

				log.Printf("[%v] read failed: %s", localAddr, rerr)
				break OuterLoop
			}
			if cc != *sz {
				continue
			}
			salt1 := string(rbuf[:8])
			log.Printf("[%s]<%s> recv %d bytes", localAddr, salt1, *sz)
			if salt == salt1 {
				lost = false
			}
		}
		if lost {
			log.Printf("[%v]<%s> lost", localAddr, salt)
		}
	}
}
