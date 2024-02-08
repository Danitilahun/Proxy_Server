package main

import (
	"bufio"    // Import bufio package for buffered I/O
	"io"       // Import io package for I/O operations
	"log"      // Import log package for logging
	"net"      // Import net package for network operations
	"net/http" // Import net/http package for HTTP protocol handling
	"strconv"  // Import strconv package for string conversions
	"sync"     // Import sync package for synchronization
	"time"     // Import time package for time-related operations
)

// Define a struct type Backend representing a backend server connection
type Backend struct {
	net.Conn               // Embed net.Conn interface for connection
	Reader   *bufio.Reader // Reader for efficient reading
	Writer   *bufio.Writer // Writer for efficient writing
}

var backendQueue chan *Backend    // Channel for managing backend connections
var requestBytes map[string]int64 // Map for storing statistics of bytes transferred for each URL path
var requestLock sync.Mutex        // Mutex for synchronization

// Initialize global variables
func init() {
	requestBytes = make(map[string]int64)  // Initialize requestBytes map
	backendQueue = make(chan *Backend, 10) // Initialize backendQueue with buffer size 10
}

// Function to retrieve a backend connection
func getBackend() (*Backend, error) {
	select {
	case be := <-backendQueue: // Try to fetch a connection from backendQueue
		return be, nil
	case <-time.After(100 * time.Millisecond): // If no connection available, wait for 100 milliseconds
		be, err := net.Dial("tcp", "127.0.0.1:8081") // Dial the backend server
		if err != nil {
			return nil, err // Return error if dialing fails
		}
		return &Backend{ // Return a new Backend instance with connection, reader, and writer
			Conn:   be,
			Reader: bufio.NewReader(be),
			Writer: bufio.NewWriter(be),
		}, nil
	}
}

// Function to queue a backend connection
func queueBackend(be *Backend) {
	select {
	case backendQueue <- be: // Try to enqueue the connection into backendQueue
	case <-time.After(1 * time.Second): // If the queue is full, wait for 1 second
		be.Close() // Close the connection
	}
}

// Function to update statistics of bytes transferred for each URL path
func updateStats(req *http.Request, resp *http.Response) int64 {
	requestLock.Lock()         // Lock to prevent concurrent access
	defer requestLock.Unlock() // Unlock after function returns

	bytes := requestBytes[req.URL.Path] + resp.ContentLength // Calculate total bytes transferred
	requestBytes[req.URL.Path] = bytes                       // Update statistics for the URL path
	return bytes                                             // Return total bytes transferred
}

// Main function
func main() {
	ln, err := net.Listen("tcp", ":8080") // Listen for incoming TCP connections on port 8080
	if err != nil {
		log.Fatalf("Failed to listen: %s", err) // Log fatal error if listening fails
	}
	for { // Infinite loop to accept incoming connections
		if conn, err := ln.Accept(); err == nil { // Accept incoming connection
			go handleConnection(conn) // Spawn a goroutine to handle the connection
		}
	}
}

// Function to handle incoming client connections
func handleConnection(conn net.Conn) {
	defer conn.Close()              // Close the connection when the function returns
	reader := bufio.NewReader(conn) // Create a reader for reading client requests

	for { // Infinite loop to handle requests
		req, err := http.ReadRequest(reader) // Read HTTP request from the client
		if err != nil {                      // Check for errors while reading request
			if err != io.EOF { // Log non-EOF errors
				log.Printf("Failed to read request: %s", err)
			}
			return // Return if error occurs
		}
		be, err := getBackend() // Retrieve a backend connection
		if err != nil {
			return // Return if error occurs while getting backend connection
		}
		if err := req.Write(be.Writer); err == nil { // Forward request to backend server
			be.Writer.Flush()                     // Flush buffered writer
			if err := req.Write(be); err == nil { // Write request to backend connection
				if resp, err := http.ReadResponse(be.Reader, req); err == nil { // Read response from backend server
					bytes := updateStats(req, resp)                          // Update statistics
					resp.Header.Set("X-Bytes", strconv.FormatInt(bytes, 10)) // Set custom header for bytes transferred
					if err := resp.Write(conn); err == nil {                 // Write response back to client
						log.Printf("%s: %d", req.URL.Path, resp.StatusCode) // Log request path and response status code
					}
				}
			}
		}
		go queueBackend(be) // Queue backend connection for reuse
	}
}
