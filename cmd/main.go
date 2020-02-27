package main

import (
	"bufio"
	"bytes"
	"fmt"
	"http/pkg"
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
		go handleConn(conn)
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

		patterns := []string{"/", "/html.html", "/favicon.ico", "/text.txt", "/images.html", "/1.jpg", "/task.pdf"}

		if method != "GET" {
			return
		}
		if protocol != "HTTP/1.1" {
			return
		}

		if strings.Contains(request, "?download") {
			indexDownload := strings.IndexByte(request, '?')
			var fileName = request[:indexDownload]
			fileName = strings.TrimPrefix(fileName, "/")
			fileName = pkg.ServerFiles + fileName
			contentType := pkg.ContentDownload
			err := headerWriter(conn, fileName, contentType, request)
			if err != nil {
				log.Printf("can't sent data to header writer: %v", err)
			}
			return
		}

		for _, pattern := range patterns {
			if request == pattern {
				switch request {
				case "/":
					err := sendFile(conn, "index.html", request)
					if err != nil {
						log.Printf("can't process the request: %v", err)
					}
				case "/html.html":
					err := sendFile(conn, "html.html", request)
					if err != nil {
						log.Printf("can't process the request: %v", err)
					}
				case "favicon.ico":
					err := sendFile(conn, "favicon.ico", request)
					if err != nil {
						log.Printf("can't process the request: %v", err)
					}
				case "/text.txt":
					err := sendFile(conn, "text.txt", request)
					if err != nil {
						log.Printf("can't process the request: %v", err)
					}
				case "/images.html":
					err := sendFile(conn, "images.html", request)
					if err != nil {
						log.Printf("can't process the request: %v", err)
					}
				case "/1.jpg":
					err := sendFile(conn, "1.jpg", request)
					if err != nil {
						log.Printf("can't process the request: %v", err)
					}
				case "/task.pdf":
					err := sendFile(conn, "task.pdf", request)
					if err != nil {
						log.Printf("can't process the request: %v", err)
					}
				default:
					err := sendFile(conn, "html404.html", request)
					if err != nil {
						log.Printf("can't process the request: %v", err)
					}
				}
				return

			}/*else {
				err := sendFile(conn, "html404.html", request)
				if err != nil {
					log.Printf("can't process the request: %v", err)
				}
			}*/
		}
		return

	}

}

func sendFile(conn net.Conn, fileName, request string) (err error) {
	contentType := ""
	part := strings.IndexByte(fileName, '.')
	fileFormat := strings.TrimPrefix(fileName, fileName[:part])

	if fileFormat == ".html" {
		fileName = pkg.ServerFiles + fileName
		contentType = pkg.ContentHTML
		err = headerWriter(conn, fileName, contentType, request)
		if err != nil {
			log.Printf("can't sent data to header writer: %v", err)
		}
		return nil
	}

	if fileFormat == ".jpg" {
		fileName = pkg.ServerFiles + fileName
		contentType = pkg.ContentJPG
		err = headerWriter(conn, fileName, contentType, request)
		if err != nil {
			log.Printf("can't sent data to header writer: %v", err)
		}
		return nil
	}

	if fileFormat == ".ico" {
		fileName = pkg.ServerFiles + fileName
		contentType = pkg.ContentICO
		err = headerWriter(conn, fileName, contentType, request)
		if err != nil {
			log.Printf("can't sent data to header writer: %v", err)
		}
		return nil
	}

	if fileFormat == ".txt" {
		fileName = pkg.ServerFiles + fileName
		contentType = pkg.ContentTEXT
		err = headerWriter(conn, fileName, contentType, request)
		if err != nil {
			log.Printf("can't sent data to header writer: %v", err)
		}
		return nil
	}

	if fileFormat == ".pdf" {
		fileName = pkg.ServerFiles + fileName
		contentType = pkg.ContentPDF
		err = headerWriter(conn, fileName, contentType, request)
		if err != nil {
			log.Printf("can't sent data to header writer: %v", err)
		}
		return nil
	}

	return nil
}

func headerWriter(conn net.Conn, fileName, contentType, request string) (err error) {
	byteBuff := bytes.Buffer{}
	writer := bufio.NewWriter(conn)
	bytesFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("can't open file: %s, error: %v", fileName, err)
		return err
	}

	log.Printf("request: %s", request)
	_, err = byteBuff.WriteString("HTTP/1.1 200 OK\r\n")
	if err != nil {
		log.Printf("can't write to buffer: %v", err)
	}

	_, err = byteBuff.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(bytesFile)))
	if err != nil {
		log.Printf("can't write to buffer: %v", err)
	}

	_, err = byteBuff.WriteString("Content-Type: " + contentType + "\r\n")
	if err != nil {
		log.Printf("can't write to buffer: %v", err)
	}

	_, err = byteBuff.WriteString("Connection: Close\r\n")
	if err != nil {
		log.Printf("can't write to buffer: %v", err)
	}

	_, err = byteBuff.WriteString("\r\n")
	if err != nil {
		log.Printf("can't write to buffer: %v", err)
	}

	_, err = byteBuff.Write(bytesFile)
	if err != nil {
		log.Printf("can't write to buffer: %v", err)
	}

	_, err = byteBuff.WriteTo(writer)
	if err != nil {
		log.Printf("error write from buffer to writer: %v", err)
	}

	err = writer.Flush()
	if err != nil {
		log.Printf("error to response: %v", err)
	}

	log.Printf("response on: %s", request)
	return nil
}
