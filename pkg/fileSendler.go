package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

func SendFile(conn net.Conn, fileName, request string) (err error) {
	contentType := ""
	part := strings.IndexByte(fileName, '.')
	fileFormat := strings.TrimPrefix(fileName, fileName[:part])

	if fileFormat == ".html" {
		fileName = ServerFiles + fileName
		contentType = ContentHTML
		err = HeaderWriter(conn, fileName, contentType, request)
		if err != nil {
			log.Printf("can't sent data to header writer: %v", err)
		}
		return nil
	}

	if fileFormat == ".jpg" {
		fileName = ServerFiles + fileName
		contentType = ContentJPG
		err = HeaderWriter(conn, fileName, contentType, request)
		if err != nil {
			log.Printf("can't sent data to header writer: %v", err)
		}
		return nil
	}

	if fileFormat == ".png" {
		fileName = ServerFiles + fileName
		contentType = ContentJPG
		err = HeaderWriter(conn, fileName, contentType, request)
		if err != nil {
			log.Printf("can't sent data to header writer: %v", err)
		}
		return nil
	}

	if fileFormat == ".txt" {
		fileName = ServerFiles + fileName
		contentType = ContentTEXT
		err = HeaderWriter(conn, fileName, contentType, request)
		if err != nil {
			log.Printf("can't sent data to header writer: %v", err)
		}
		return nil
	}

	if fileFormat == ".pdf" {
		fileName = ServerFiles + fileName
		contentType = ContentPDF
		err = HeaderWriter(conn, fileName, contentType, request)
		if err != nil {
			log.Printf("can't sent data to header writer: %v", err)
		}
		return nil
	}

	return nil
}

func HeaderWriter(conn net.Conn, fileName, contentType, request string) (err error) {
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

