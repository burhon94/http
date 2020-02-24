package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	host := "0.0.0.0"
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "9999"
	}
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("start server on: %s", addr)
	err := start(addr)
	if err != nil {
		log.Fatalf("can't start server on: %s, error: %v", addr, err)
	}
}

func start(addr string) (err error) {
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
		handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("can't close handle connect: %v", err)
			return
		}
	}()
	log.Print("read client request")
	reader := bufio.NewReaderSize(conn, 4096)
	writer := bufio.NewWriter(conn)
	counter := 0
	buf := [4096]byte{}
	for {
		if counter == 4096 {
			log.Printf("too long request header")
			_, _ = writer.WriteString("HTTP/1.1 413 Payload Too Large\r\n")
			_, _ = writer.WriteString("Content-Length: 0\r\n")
			_, _ = writer.WriteString("Connection: close\r\n")
			_, _ = writer.WriteString("\r\n")
			err := writer.Flush()
			if err != nil {
				log.Printf("can't sent response: %v", err)
			}
			return
		}
		read, err := reader.ReadByte()
		if err != nil {
			log.Printf("can't read request line: %v", err)
			_, _ = writer.WriteString("HTTP/1.1 400 Bad Request\r\n")
			_, _ = writer.WriteString("Content-Length: 0\r\n")
			_, _ = writer.WriteString("Connection: close\r\n")
			_, _ = writer.WriteString("\r\n")
			err := writer.Flush()
			if err != nil {
				log.Printf("can't sent response: %v", err)
			}
			return
		}
		buf[counter] = read
		counter++
		if counter < 4 {
			continue
		}
		if string(buf[counter-4:counter]) == "\r\n\r\n" {
			break
		}
	}
	headersStr := string(buf[:counter-4])
	requestHeaderParts := strings.Split(headersStr, "\r\n")
	log.Print("parse request line")
	requestLine := requestHeaderParts[0]
	requestParts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(requestParts) != 3 {
		return
	}
	method, request, protocol := requestParts[0], requestParts[1], requestParts[2]

	for {
		b := request != "/"
		b2 := request != "/html.html"
		ic := request != "/favicon.ico"
		b3 := request != "/text.txt"
		b4 := request != "/images.html"
		im := request != "/img/1.jpg"
		b5 := request != "/task.pdf"
		b6 := request != "/task.pdf?download"
		if b && ic && b2 && b3 && b4 && im && b5 && b6 {

			html404 := `<!doctype html>
			<html lang="en">
			<meta charset="UTF-8">
			<meta name="viewport"
			content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">

			<title>The page is not found</title>

			<style type="text/css">
				body,h1,h2{margin:0}h2,h3{padding:.5em}h1,h2,h3{text-align:center;color:#fff}body{background-color:#fff;color:#000;font-size:.9em;font-family:sans-serif,helvetica;padding:0}:link,:visited{color:#c00}a:hover{color:#f50}h1{padding:.6em 2em .4em;background-color:#294172;font-weight:400;font-size:1.75em;border-bottom:2px solid #000}h1 strong,h2{font-weight:700}h1 strong{font-size:1.5em}h2{background-color:#3C6EB4;font-size:1.1em;border-bottom:2px solid #294172}h3{background-color:red}hr{display:none}.content{padding:1em 5em}.alert{border:2px solid #000}img{border:2px solid #fff;padding:2px;margin:2px}a:hover img{border:2px solid #294172}
			</style>
			</head>

			<body cz-shortcut-listen="true">
			<h1><strong>HTTP 404</strong></h1>

			<div class="content">

			<h3>The page you are looking for is not found.</h3>

			</div>


			</body></html>`

			log.Printf("request: %s", request)
			_, _ = writer.WriteString("HTTP/1.1 404 Page Not Found\r\n")
			_, _ = writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(html404)))
			_, _ = writer.WriteString("Content-Type: text/html\r\n")
			_, _ = writer.WriteString("Connection: Close\r\n")
			_, _ = writer.WriteString("\r\n")
			_, _ = writer.WriteString(html404)
			err := writer.Flush()
			if err != nil {
				log.Printf("can't sent response: %v", err)
			}
			log.Printf("response on: %s", request)
		}

		if method == "GET" && request == "/" && protocol == "HTTP/1.1" {
			index := `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>First</title>
</head>
<body bgcolor="#FFFFB5">
<h1>Hello, from first html</h1>
<a href="./html.html">HTML</a><br>
<a href="./text.txt">Some text</a><br>
<a href="./images.html">images</a><br>
<a href="task.pdf">pdf</a><br>
</body>
</html>`
			log.Printf("request: %s", request)
			_, _ = writer.WriteString("HTTP/1.1 200 OK\r\n")
			_, _ = writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(index)))
			_, _ = writer.WriteString("Content-Type: text/html\r\n")
			_, _ = writer.WriteString("Connection: Close\r\n")
			_, _ = writer.WriteString("\r\n")
			_, _ = writer.WriteString(index)
			err := writer.Flush()
			if err != nil {
				log.Printf("can't sent response: %v", err)
			}
			log.Printf("response on: %s", request)
		}

		if method == "GET" && request == "/favicon.ico" && protocol == "HTTP/1.1" {
			bytes, err := ioutil.ReadFile("./server/img/icon.png")
			if err != nil {
				log.Printf("can't open file")
			}
			log.Printf("request: %s", request)
			_, _ = writer.WriteString("HTTP/1.1 200 OK\r\n")
			_, _ = writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(bytes)))
			_, _ = writer.WriteString("Content-Type: image/png\r\n")
			_, _ = writer.WriteString("Connection: Close\r\n")
			_, _ = writer.WriteString("\r\n")
			_, _ = writer.Write(bytes)
			err = writer.Flush()
			if err != nil {
				log.Printf("can't sent response: %v", err)
			}
			log.Printf("response on: %s", request)
		}

		if method == "GET" && request == "/html.html" && protocol == "HTTP/1.1" {
			log.Printf("request: %s", request)
			bytes, err := ioutil.ReadFile("./server/html.html")
			if err != nil {
				log.Printf("can't load html.html: %v", err)
			}
			_, _ = writer.WriteString("HTTP/1.1 200 OK\r\n")
			_, _ = writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(bytes)))
			_, _ = writer.WriteString("Content-Type: text/html\r\n")
			_, _ = writer.WriteString("Connection: Close\r\n")
			_, _ = writer.WriteString("\r\n")
			_, _ = writer.Write(bytes)
			err = writer.Flush()
			if err != nil {
				log.Printf("can't sent response: %v", err)
			}
			log.Printf("response on: %s", request)
		}

		if method == "GET" && request == "/text.txt" && protocol == "HTTP/1.1" {
			log.Printf("request: %s", request)
			bytes, err := ioutil.ReadFile("./server/someText.txt")
			if err != nil {
				log.Printf("can't load someText.txt: %v", err)
			}
			_, _ = writer.WriteString("HTTP/1.1 200 OK\r\n")
			_, _ = writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(bytes)))
			_, _ = writer.WriteString("Content-Type: text/plain\r\n")
			_, _ = writer.WriteString("Connection: Close\r\n")
			_, _ = writer.WriteString("\r\n")
			_, _ = writer.Write(bytes)
			err = writer.Flush()
			if err != nil {
				log.Printf("can't sent response: %v", err)
			}
			log.Printf("response on: %s", request)
		}

		if method == "GET" && request == "/images.html" && protocol == "HTTP/1.1" {
			log.Printf("request: %s", request)
			bytes, err := ioutil.ReadFile("./server/images.html")
			if err != nil {
				log.Printf("can't load images.html: %v", err)
			}
			_, _ = writer.WriteString("HTTP/1.1 200 OK\r\n")
			_, _ = writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(bytes)))
			_, _ = writer.WriteString("Connection: Close\r\n")
			_, _ = writer.WriteString("Content-type: text/html\r\n")
			_, _ = writer.WriteString("\r\n")
			_, _ = writer.Write(bytes)
			err = writer.Flush()
			if err != nil {
				log.Printf("can't sent response: %v", err)
			}
			log.Printf("response on: %s", request)
		}

		if method == "GET" && request == "/img/1.jpg" && protocol == "HTTP/1.1" {
			log.Printf("request: %s", request)
			bytesIMG, err := ioutil.ReadFile("./server/img/1.jpg")
			if err != nil {
				log.Printf("can't load 1.jpg: %v", err)
			}
			_, _ = writer.WriteString("HTTP/1.1 200 OK\r\n")
			_, _ = writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(bytesIMG)))
			_, _ = writer.WriteString("Connection: Close\r\n")
			_, _ = writer.WriteString("Content-type: media\r\n")
			_, _ = writer.WriteString("\r\n")
			_, _ = writer.Write(bytesIMG)
			err = writer.Flush()
			if err != nil {
				log.Printf("can't sent response: %v", err)
			}
			log.Printf("response on: %s", request)
		}

		if method == "GET" && request == "/task.pdf" && protocol == "HTTP/1.1" {
			log.Printf("request: %s", request)
			bytesPDF, err := ioutil.ReadFile("./server/file/html.pdf")
			if err != nil {
				log.Printf("can't load task.pdf: %v", err)
			}
			_, _ = writer.WriteString("HTTP/1.1 200 OK\r\n")
			_, _ = writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(bytesPDF)))
			_, _ = writer.WriteString("Connection: Close\r\n")
			_, _ = writer.WriteString("Content-type: application/pdf\r\n")
			_, _ = writer.WriteString("\r\n")
			_, _ = writer.Write(bytesPDF)
			err = writer.Flush()
			if err != nil {
				log.Printf("can't sent response: %v", err)
			}
			log.Printf("response on: %s", request)
		}

		if method == "GET" && request == "/task.pdf?download" && protocol == "HTTP/1.1" {
			log.Printf("request: %s", request)
			bytesPDF, err := ioutil.ReadFile("./server/file/html.pdf")
			if err != nil {
				log.Printf("can't load task.pdf: %v", err)
			}
			_, _ = writer.WriteString("HTTP/1.1 200 OK\r\n")
			_, _ = writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(bytesPDF)))
			_, _ = writer.WriteString("Connection: Close\r\n")
			_, _ = writer.WriteString("Content-type: application/download\r\n")
			_, _ = writer.WriteString("\r\n")
			_, _ = writer.Write(bytesPDF)
			err = writer.Flush()
			if err != nil {
				log.Printf("can't sent response: %v", err)
			}
			log.Printf("response on: %s", request)
		}

		return

	}

}
