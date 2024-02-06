# HTTP Reverse Proxy in Go

This Go program acts as a simple HTTP reverse proxy. It listens for incoming TCP connections on port 8080, reads HTTP requests from clients, forwards them to another server (localhost:8081), reads the response from that server, and then writes the response back to the client.

## Implementation

The program consists of two main functions: `main()` and `handleConnection()`.

### `main()`

- It listens for incoming TCP connections on port 8080 using `net.Listen()`.
- It accepts incoming connections using `ln.Accept()` in an infinite loop and launches a goroutine to handle each connection.

### `handleConnection(conn net.Conn)`

- This function is responsible for handling each incoming connection.
- It reads HTTP requests from the client connection using `bufio.NewReader()`.
- For each request received, it dials a connection to the backend server (localhost:8081) using `net.Dial()`.
- It forwards the HTTP request to the backend server and reads the response using `http.ReadResponse()`.
- Finally, it writes the response back to the client connection using `resp.Write()`.

## Error Handling

- The program handles errors such as failing to listen for connections, failing to read requests, and failing to write responses. Error messages are logged using `log.Printf()`.

## Note

- This program provides a basic example of a reverse proxy and may require additional features such as request/response modification, load balancing, and security measures for production use.
