# Backend Proxy with Connection Pooling

## Overview

This Go program implements a TCP proxy server with connection pooling for backend servers. It listens for incoming connections on port 8080, forwards HTTP requests to a backend server, and tracks statistics of bytes transferred for each URL path.

## Components

- **Backend Struct**: Represents a backend server connection with a TCP connection, bufio.Reader, and bufio.Writer for efficient I/O operations.

- **Global Variables**:

  - `backendQueue`: A buffered channel for managing backend connections.
  - `requestBytes`: A map storing statistics of bytes transferred for each URL path.
  - `requestLock`: A mutex for synchronizing access to `requestBytes`.

- **Initialization**:

  - The `init()` function initializes `requestBytes` map and `backendQueue` channel with a buffer size of 10.

- **getBackend() Function**:

  - Retrieves a backend connection from `backendQueue` or establishes a new connection if the queue is empty.
  - Uses a timeout of 100 milliseconds for connection retrieval.

- **queueBackend() Function**:

  - Queues a backend connection back into `backendQueue` with a timeout of 1 second.

- **updateStats() Function**:

  - Updates statistics of bytes transferred for each URL path.
  - Locked using `requestLock` to prevent concurrent access.

- **Main Function**:

  - Listens for incoming TCP connections on port 8080.
  - Spawns a goroutine to handle each connection.

- **handleConnection() Function**:
  - Handles each client connection.
  - Reads an HTTP request from the client and retrieves a backend connection using `getBackend()`.
  - Forwards the request to the backend server, reads the response, updates statistics, modifies the response with a custom header, and sends it back to the client.
  - Queues the backend connection back into `backendQueue` asynchronously.

## Usage

- Run the program by compiling and executing the Go file.
- The program will start listening for incoming TCP connections on port 8080.
- Clients can connect to the proxy server and send HTTP requests, which will be forwarded to the backend server.
- Statistics of bytes transferred for each URL path will be tracked.
