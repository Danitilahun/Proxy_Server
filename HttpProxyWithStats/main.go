package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

// Declare global variables for tracking request bytes and a mutex for synchronization
var requestBytes map[string]int64
var requestLock sync.Mutex

// Initialize the requestBytes map
func init() {
	requestBytes = make(map[string]int64)
}

// Update request bytes and return the total bytes transferred for a specific URL path
func updateStats(req *http.Request, resp *http.Response) int64 {
	// Lock to prevent concurrent access to requestBytes map
	requestLock.Lock()
	defer requestLock.Unlock()

	// Calculate total bytes transferred and update requestBytes map
	bytes := requestBytes[req.URL.Path] + resp.ContentLength
	requestBytes[req.URL.Path] = bytes
	return bytes
}

func main() {
	// Listen for incoming TCP connections on port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}
	for {
		// Accept incoming connections
		if conn, err := ln.Accept(); err == nil {
			// Handle each connection concurrently
			go handleConnection(conn)
		}
	}
}

// Handle incoming client connections
func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// Read HTTP request from client
		req, err := http.ReadRequest(reader)
		if err != nil {
			// Handle errors while reading request
			if err != io.EOF {
				log.Printf("Failed to read request: %s", err)
			}
			return
		}
		// Dial backend server
		if be, err := net.Dial("tcp", "127.0.0.1:8081"); err == nil {
			be_reader := bufio.NewReader(be)
			// Forward request to backend server
			if err := req.Write(be); err == nil {
				// Read response from backend server
				if resp, err := http.ReadResponse(be_reader, req); err == nil {
					// Update statistics and add custom header
					bytes := updateStats(req, resp)
					resp.Header.Set("X-Bytes", strconv.FormatInt(bytes, 10))
					// Write response back to client
					if err := resp.Write(conn); err == nil {
						log.Printf("%s: %d", req.URL.Path, resp.StatusCode)
					}
				}
			}
		}
	}
}
