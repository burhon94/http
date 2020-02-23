package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

func Test_ServerConnOK(t *testing.T) {
	addr := "0.0.0.0:9988"
	go func() {
		log.Printf("start server")
		err:= start(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
	}()
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't dial: %v", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			t.Fatalf("can't close dial: %v", err)
		}
	}()
}

func Test_ServerIndex(t *testing.T) {
	addr := "0.0.0.0:9922"
	go func() {
		log.Printf("start server")
		err:= start(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
	}()
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't dial: %v", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			t.Fatalf("can't close dial: %v", err)
		}
	}()
	write := bufio.NewWriter(conn)
	write.WriteString("GET / HTTP/1.1\r\n")
	write.WriteString("Host: localhost\r\n")
	write.WriteString("\r\n")
	write.Flush()
	bytes, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Fatalf("can't read response from server: %v", err)
	}
	response := string(bytes)
	if !strings.Contains(response, "200 OK") {
		t.Fatalf("just it be 200 OK: %s", response)
	}
}

func Test_ServerImages(t *testing.T) {
	addr := "0.0.0.0:8870"
	go func() {
		log.Printf("start server")
		err:= start(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
	}()
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't dial: %v", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			t.Fatalf("can't close dial: %v", err)
		}
	}()
	write := bufio.NewWriter(conn)
	write.WriteString("GET /images.html HTTP/1.1\r\n")
	write.WriteString("Host: localhost\r\n")
	write.WriteString("\r\n")
	write.Flush()
	bytes, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Fatalf("can't read response from server: %v", err)
	}
	response := string(bytes)
	if !strings.Contains(response, "200 OK") {
		t.Fatalf("just it be 200 OK: %s", response)
	}
}