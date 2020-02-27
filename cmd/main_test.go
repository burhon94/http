package cmd

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
		log.Printf("Start server")
		err:= Start(addr)
		if err != nil {
			t.Fatalf("can't Start server: %v", err)
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
		log.Printf("Start server")
		err:= Start(addr)
		if err != nil {
			t.Fatalf("can't Start server: %v", err)
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

func Test_ServerHTML(t *testing.T) {
	addr := "0.0.0.0:9933"
	go func() {
		log.Printf("Start server")
		err:= Start(addr)
		if err != nil {
			t.Fatalf("can't Start server: %v", err)
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
	write.WriteString("GET /html.html HTTP/1.1\r\n")
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

func Test_ServerTEXT(t *testing.T) {
	addr := "0.0.0.0:9977"
	go func() {
		log.Printf("Start server")
		err:= Start(addr)
		if err != nil {
			t.Fatalf("can't Start server: %v", err)
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
	write.WriteString("GET /text.txt HTTP/1.1\r\n")
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

func Test_ServerImageHTML(t *testing.T) {
	addr := "0.0.0.0:8870"
	go func() {
		log.Printf("Start server")
		err:= Start(addr)
		if err != nil {
			t.Fatalf("can't Start server: %v", err)
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

func Test_ServerImages(t *testing.T) {
	addr := "0.0.0.0:8892"
	go func() {
		log.Printf("Start server")
		err:= Start(addr)
		if err != nil {
			t.Fatalf("can't Start server: %v", err)
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
	write.WriteString("GET /img/1.jpg HTTP/1.1\r\n")
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

func Test_ServerPDF(t *testing.T) {
	addr := "0.0.0.0:8852"
	go func() {
		log.Printf("Start server")
		err:= Start(addr)
		if err != nil {
			t.Fatalf("can't Start server: %v", err)
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
	write.WriteString("GET /task.pdf HTTP/1.1\r\n")
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

func Test_Server404(t *testing.T) {
	addr := "0.0.0.0:9966"
	go func() {
		log.Printf("Start server")
		err:= Start(addr)
		if err != nil {
			t.Fatalf("can't Start server: %v", err)
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
	write.WriteString("GET /404.html HTTP/1.1\r\n")
	write.WriteString("Host: localhost\r\n")
	write.WriteString("\r\n")
	write.Flush()
	bytes, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Fatalf("can't read response from server: %v", err)
	}
	response := string(bytes)
	if !strings.Contains(response, "404 Page Not Found") {
		t.Fatalf("response just 404 Page Not Found: %s", response)
	}
}