package pkg

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"sync"
)

func HTTP400(conn net.Conn) error {
	writer := bufio.NewWriter(conn)
	_, _ = writer.WriteString("HTTP/1.1 400 Bad Request\r\n")
	_, _ = writer.WriteString("Content-Length: 0\r\n")
	_, _ = writer.WriteString("Connection: close\r\n")
	_, _ = writer.WriteString("\r\n")
	err := writer.Flush()
	if err != nil {
		log.Printf("can't sent response: %v", err)
	}
	return err
}

func HandleConn(conn net.Conn) (err error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("can't close handle connect: %v", err)
			return
		}
	}()

	log.Print("read client request")
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("can't read request line: %v", err)
		err := HTTP400(conn)
		if err != nil {
			log.Printf("can't sent response: %v ", err)
		}
		return err
	}

	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) != 3 {
		log.Printf("can't read request line: %v", err)
		err := HTTP400(conn)
		if err != nil {
			log.Printf("can't sent response: %v ", err)
		}
		return err
	}

	requestParts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(requestParts) != 3 {
		return
	}

	wg := sync.WaitGroup{}
	method, request, protocol := requestParts[0], requestParts[1], requestParts[2]
	for {
		wg.Add(1)
		err := handleRequest(conn, method, request, protocol, wg)
		if err != nil {
			log.Printf("error while proccessed request line: %v", err)
		}
		wg.Wait()

		return nil
	}
}

func handleRequest(conn net.Conn, method string, request string, protocol string, wg sync.WaitGroup) (err error) {
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
		fileName = ServerFiles + fileName
		contentType := ContentDownload
		err := HeaderWriter(conn, fileName, contentType, request)
		if err != nil {
			log.Printf("can't sent data to header writer: %v", err)
			return err
		}
		return nil
	}
	if request == "/" {
		err := SendFile(conn, "index.html", request)
		if err != nil {
			log.Printf("can't process the request: %v", err)
		}
		return err
	}
	if request == "/favicon.ico" {
		err := SendFile(conn, "icon.png", request)
		if err != nil {
			log.Printf("can't process the request: %v", err)
		}
		return err
	}
	var fileName = request[:]
	fileName = strings.TrimPrefix(fileName, "/")
	serverFiles, err := ioutil.ReadDir(ServerFiles)
	if err != nil {
		log.Printf("can't check server files: %s, error %v", ServerFiles, err)
	}
	for _, serverFile := range serverFiles {
		if fileName == serverFile.Name() {
			err := SendFile(conn, fileName, request)
			if err != nil {
				log.Printf("can't process the request: %v", err)
				return err
			}
			return nil
		}
	}
	err = SendFile(conn, "html404.html", request)
	if err != nil {
		log.Printf("can't process the request: %v", err)
	}
	wg.Done()
	return nil
}
