package main

import (
	"flag"
	"log"
	"net"
	"peerchat/internal/client"
)

func main() {
	addr := flag.String("addr", ":8080", "client address")
	tcpAddr, err := net.ResolveTCPAddr("tcp", *addr)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := client.NewClient(tcpAddr)
	defer client.Close()
	client.Start()
}
