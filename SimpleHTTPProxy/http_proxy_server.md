# HTTP Proxy Server in Go

This Go code implements a basic HTTP proxy server that listens on port 8080. Here's a summary of what it does:

1. It listens for incoming TCP connections on port 8080.
2. For each incoming connection:
   - It reads the HTTP request sent by the client.
   - It establishes a connection to a backend server at `127.0.0.1:8081`.
   - It forwards the received HTTP request to the backend server.
   - It reads the response from the backend server.
   - It sends the received response back to the client.
   - It logs the path and status code of the request-response pair.
3. It closes the client connection after completing the request-response cycle.

This code effectively acts as a simple HTTP proxy server, forwarding requests from clients to a backend server and relaying the responses back to the clients.
