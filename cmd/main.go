package main

import (
	"fmt"
	"http/pkg"
	"log"
	"net"
	"os"
)

func main() {
	host := "0.0.0.0"
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "9999"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Start server on: %s", addr)

	err := Start(addr)
	if err != nil {
		log.Fatalf("can't Start server on: %s, error: %v", addr, err)
	}
}

func Start(addr string) (err error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("can't listen on: %s, error: %v", addr, err)
		return err
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			log.Fatalf("can't close server listener: %v", err)
			return
		}
	}()
	for {
		conn, err := listener.Accept()
		log.Print("try accept connection")
		if err != nil {
			log.Printf("can't accept, error: %v", err)
			continue
		}
		log.Print("accept success")
		go func() {
			err := pkg.HandleConn(conn)
			if err != nil {
				log.Printf("can't handler connect: %v", err)
			}
		}()
	}
}

