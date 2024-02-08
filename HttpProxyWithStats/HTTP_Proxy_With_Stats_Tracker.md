# HTTP Proxy with Stats Tracker

## Overview

This Go program serves as an HTTP proxy server that listens on port 8080 for incoming TCP connections. It forwards incoming HTTP requests to a backend server and tracks statistics related to the bytes transferred for each URL path.

## Components

- **requestBytes**: A map variable that stores statistics about the number of bytes transferred for each URL path in HTTP requests.

- **requestLock**: A mutex variable used for synchronizing access to the `requestBytes` map to prevent data race conditions.

- **init()**: Initializes the `requestBytes` map.

- **updateStats()**: Updates statistics related to bytes transferred for each URL path based on the HTTP request and response.

- **main()**: The entry point of the program. It listens for incoming TCP connections on port 8080 and spawns a goroutine to handle each connection.

- **handleConnection()**: Handles each client connection. It reads HTTP requests from clients, forwards them to a backend server, reads the responses from the backend server, updates statistics, and sends the responses back to clients.

## Functionality

1. The program listens for incoming TCP connections on port 8080 and accepts connections from clients.

2. For each client connection, it reads an HTTP request from the client.

3. It establishes a connection to a backend server (localhost:8081) and forwards the HTTP request to the backend server.

4. Upon receiving the response from the backend server, it updates statistics about the bytes transferred for the specific URL path in the request.

5. It adds a custom header `X-Bytes` to the response containing the total bytes transferred for that URL path.

6. Finally, it sends the response back to the client.

## Usage

- The program is designed to be run as a standalone application.

- It can be compiled and executed on a system with Go installed.

- It requires a backend server running on localhost:8081 to forward HTTP requests.

- To compile and run the program:
