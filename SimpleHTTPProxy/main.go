package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
)

func main() {
	// Listen for incoming TCP connections on port 8080
	if listener, err := net.Listen("tcp", ":8080"); err == nil {
		for {
			// Accept incoming connection
			if clientConn, err := listener.Accept(); err == nil {
				// Create a buffered reader for the client connection
				clientReader := bufio.NewReader(clientConn)

				// Read the HTTP request
				if clientReq, err := http.ReadRequest(clientReader); err == nil {
					// Dial to backend server
					if backendConn, err := net.Dial("tcp", "127.0.0.1:8081"); err == nil {
						// Create a buffered reader for the backend connection
						backendReader := bufio.NewReader(backendConn)

						// Forward the HTTP request to backend server
						if err := clientReq.Write(backendConn); err == nil {
							// Read the response from backend server
							if backendResp, err := http.ReadResponse(backendReader, clientReq); err == nil {
								// Close the response body after sending
								backendResp.Close = true

								// Send the response back to the client
								if err := backendResp.Write(clientConn); err == nil {
									// Log the request path and status code
									log.Printf("%s: %d", clientReq.URL.Path, backendResp.StatusCode)
								}

								// Close the client connection
								clientConn.Close()
							}
						}
					}
				}
			}
		}
	}
}
