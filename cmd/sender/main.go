package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/atotto/clipboard"
)

func main() {
	targetIP := flag.String("ip", "", "Target Windows IP address (Required)")
	targetPort := flag.String("port", "9999", "Target Windows Port")
	flag.Parse()

	if *targetIP == "" {
		fmt.Println("Error: Please provide the Windows IP address using -ip")
		flag.Usage()
		return
	}

	targetAddr := fmt.Sprintf("%s:%s", *targetIP, *targetPort)
	log.Printf("Sender started. Monitoring clipboard...")
	log.Printf("Target Receiver: %s", targetAddr)

	var lastContent string
	// Initial read to avoid sending existing content on startup (optional, but good for UX)
	// Or we can start empty. Let's start empty to sync current state if it changes or just wait.
	// To be safe, let's just initialize lastContent with current to avoid resending immediately on boot
	// unless the user copied something *new*.
	// However, if I restart the sender, I might want to sync the current state.
	// Let's start with lastContent = "" so the first check will sync if there is content.
	// But if the clipboard has "A", and we start, we read "A". "A" != "". Syncs "A".
	// This is usually desired behavior (sync current state on connect).

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		content, err := clipboard.ReadAll()
		if err != nil {
			log.Printf("Error reading clipboard: %v", err)
			continue
		}

		if content != "" && content != lastContent {
			log.Printf("Clipboard changed. Length: %d. Sending...", len(content))
			
			if err := sendData(targetAddr, content); err != nil {
				log.Printf("Failed to sync: %v", err)
				// Note: We do NOT update lastContent if send fails, 
				// so we retry on next tick (because content != lastContent still holds).
				// However, if the error persists, we spam logs. 
				// To avoid log spam, we could update lastContent, or use a backoff.
				// For this simple tool, not updating means we retry, which is "robust" but noisy.
				// Let's update lastContent to prevent infinite retry loops on network down,
				// or keep it to ensure eventual consistency?
				// "Reliability" suggests we should try to sync eventually.
				// But "Polling" will re-trigger every 500ms.
				// Let's NOT update lastContent, so it retries.
			} else {
				lastContent = content
				log.Println("Sync successful.")
			}
		}
	}
}

func sendData(address, content string) error {
	// Connect to the server
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	defer conn.Close()

	// Prepare data
	data := []byte(content)
	length := uint32(len(data))
	
	// Protocol: 4 bytes length (BigEndian) + Content
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, length)

	// Write Header
	if _, err := conn.Write(header); err != nil {
		return fmt.Errorf("write header error: %w", err)
	}

	// Write Body
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("write body error: %w", err)
	}

	return nil
}
