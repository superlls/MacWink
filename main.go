package main

import (
	"context"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
)

func main() {
	peerFlag := flag.String("peer", "", "Peer IP/host (required). You may also pass host:port.")
	portFlag := flag.String("port", "9999", "Local listen port")
	intervalFlag := flag.Duration("interval", 500*time.Millisecond, "Clipboard poll interval")
	flag.Parse()

	if *peerFlag == "" {
		fmt.Fprintln(os.Stderr, "Error: -peer is required, e.g. go run main.go -peer 192.168.1.100")
		flag.Usage()
		os.Exit(2)
	}

	peerAddr, err := normalizePeerAddr(*peerFlag, *portFlag)
	if err != nil {
		log.Fatalf("Invalid -peer value: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tracking := &trackingValue{}
	outbox := make(chan string, 1)

	go func() {
		if err := runServer(ctx, ":"+*portFlag, tracking); err != nil {
			log.Printf("server stopped: %v", err)
			stop()
		}
	}()
	go runClient(ctx, peerAddr, outbox)
	go runWatcher(ctx, *intervalFlag, tracking, outbox)

	log.Printf("Peer started. Listening on :%s, peer=%s", *portFlag, peerAddr)
	<-ctx.Done()
}

type trackingValue struct {
	mu   sync.RWMutex
	last string
}

func (t *trackingValue) Get() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.last
}

func (t *trackingValue) Set(v string) {
	t.mu.Lock()
	t.last = v
	t.mu.Unlock()
}

func normalizePeerAddr(peer, defaultPort string) (string, error) {
	if peer == "" {
		return "", errors.New("empty peer")
	}
	if _, _, err := net.SplitHostPort(peer); err == nil {
		return peer, nil
	}
	return net.JoinHostPort(peer, defaultPort), nil
}

func runServer(ctx context.Context, listenAddr string, tracking *trackingValue) error {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) || ctx.Err() != nil {
				return nil
			}
			log.Printf("accept error: %v", err)
			continue
		}
		go handleIncomingConn(ctx, conn, tracking)
	}
}

const maxFrameSize = 10 * 1024 * 1024 // 10 MiB

func handleIncomingConn(ctx context.Context, conn net.Conn, tracking *trackingValue) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("incoming connection: %s", remoteAddr)

	for {
		if ctx.Err() != nil {
			return
		}
		content, err := readFrame(conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("connection closed: %s", remoteAddr)
			} else {
				log.Printf("read error from %s: %v", remoteAddr, err)
			}
			return
		}
		if content == "" {
			continue
		}

		prev := tracking.Get()
		tracking.Set(content)
		if err := clipboard.WriteAll(content); err != nil {
			tracking.Set(prev)
			log.Printf("clipboard write error: %v", err)
			continue
		}
		log.Printf("received %d bytes from %s", len(content), remoteAddr)
	}
}

func runWatcher(ctx context.Context, interval time.Duration, tracking *trackingValue, outbox chan string) {
	lastLocal := ""
	if initial, err := clipboard.ReadAll(); err == nil {
		lastLocal = initial
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			content, err := clipboard.ReadAll()
			if err != nil {
				log.Printf("clipboard read error: %v", err)
				continue
			}
			if content == "" || content == lastLocal {
				continue
			}
			lastLocal = content

			if content == tracking.Get() {
				continue
			}

			sendLatest(outbox, content)
		}
	}
}

func sendLatest(outbox chan string, content string) {
	select {
	case outbox <- content:
		log.Printf("clipboard changed: %d bytes (queued)", len(content))
	default:
		select {
		case <-outbox:
		default:
		}
		outbox <- content
		log.Printf("clipboard changed: %d bytes (replaced)", len(content))
	}
}

func runClient(ctx context.Context, peerAddr string, outbox <-chan string) {
	var conn net.Conn
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	var pending string
	var hasPending bool

	backoff := 500 * time.Millisecond
	keepalive := time.NewTicker(10 * time.Second)
	defer keepalive.Stop()

	dial := func() net.Conn {
		dialer := &net.Dialer{Timeout: 3 * time.Second, KeepAlive: 30 * time.Second}
		c, err := dialer.DialContext(ctx, "tcp", peerAddr)
		if err != nil {
			log.Printf("connect to %s failed: %v (retry in %s)", peerAddr, err, backoff)
			if backoff < 8*time.Second {
				backoff *= 2
			}
			return nil
		}
		backoff = 500 * time.Millisecond
		log.Printf("connected to peer: %s", peerAddr)
		return c
	}

	for {
		if ctx.Err() != nil {
			return
		}
		if conn == nil {
			conn = dial()
			if conn == nil {
				select {
				case <-ctx.Done():
					return
				case content := <-outbox:
					pending = content
					hasPending = true
				case <-time.After(backoff):
				}
			}
			continue
		}

		for conn != nil && hasPending {
			if err := writeFrame(conn, pending); err != nil {
				log.Printf("send failed (reconnect): %v", err)
				_ = conn.Close()
				conn = nil
				break
			}
			hasPending = false
			log.Printf("sent %d bytes to %s", len(pending), peerAddr)
		}
		if conn == nil {
			continue
		}

		select {
		case <-ctx.Done():
			return
		case content := <-outbox:
			pending = content
			hasPending = true
		case <-keepalive.C:
			if err := writeFrame(conn, ""); err != nil {
				log.Printf("keepalive failed (reconnect): %v", err)
				_ = conn.Close()
				conn = nil
			}
		}
	}
}

func readFrame(r io.Reader) (string, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return "", err
	}
	n := binary.BigEndian.Uint32(lenBuf[:])
	if n == 0 {
		return "", nil
	}
	if n > maxFrameSize {
		return "", fmt.Errorf("frame too large: %d bytes", n)
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func writeFrame(w io.Writer, content string) error {
	data := []byte(content)
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(data)))
	if _, err := w.Write(header); err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	_, err := w.Write(data)
	return err
}
