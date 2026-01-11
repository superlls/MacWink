package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/atotto/clipboard"
)

func main() {
	port := flag.String("port", "9999", "Port to listen on")
	flag.Parse()

	addr := ":" + *port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
	}
	defer listener.Close()

	log.Printf("Receiver started. Listening on %s...", addr)
	log.Println("Waiting for clipboard data from Sender...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Handle connection in a goroutine to support multiple concurrent sends (though typically sequential)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("Connection accepted from %s", remoteAddr)

	for {
		// 1. Read 4-byte length prefix
		lengthBuf := make([]byte, 4)
		_, err := io.ReadFull(conn, lengthBuf)
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection closed by %s", remoteAddr)
			} else {
				log.Printf("Error reading length from %s: %v", remoteAddr, err)
			}
			return
		}

		dataLen := binary.BigEndian.Uint32(lengthBuf)

		// 2. Read the actual content
		// Limit the reader to avoid potential OOM attacks if length is huge
		dataBuf := make([]byte, dataLen)
		_, err = io.ReadFull(conn, dataBuf)
		if err != nil {
			log.Printf("Error reading data from %s: %v", remoteAddr, err)
			return
		}

		content := string(dataBuf)
		if content == "" {
			continue // Skip empty
		}

		// 3. Write to system clipboard
		err = clipboard.WriteAll(content)
		if err != nil {
			log.Printf("Error writing to clipboard: %v", err)
		} else {
			// Truncate log for very long content
			logContent := content
			if len(logContent) > 50 {
				logContent = logContent[:50] + "..."
			}
			log.Printf("Received and synced: %q (%d bytes)", logContent, dataLen)
		}
	}
}
