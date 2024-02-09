package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type Documentation struct {
	RequestBytes map[string]int64
	Lock         sync.RWMutex
}

type GetURLRequest struct {
	URL string `json:"url"`
}

type GetURLResponse struct {
	Content string `json:"content"`
}

type GetStatsResponse struct {
	Stats map[string]int64 `json:"stats"`
}

func (d *Documentation) GetURL(request GetURLRequest, response *GetURLResponse) error {
	// Fetch the content of the URL
	resp, err := http.Get(request.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Set the content to the fetched body
	response.Content = string(body)

	// Record the request bytes
	d.Lock.Lock()
	d.RequestBytes[request.URL] = int64(len(body))
	d.Lock.Unlock()

	return nil
}

func (d *Documentation) GetStats(_ struct{}, response *GetStatsResponse) error {
	// Get a read lock to access the request bytes
	d.Lock.RLock()
	defer d.Lock.RUnlock()

	// Copy the map to the response
	response.Stats = make(map[string]int64)
	for url, bytes := range d.RequestBytes {
		response.Stats[url] = bytes
	}

	return nil
}

func main() {
	// Initialize the Documentation struct
	doc := &Documentation{
		RequestBytes: make(map[string]int64),
	}

	// Register the Documentation service
	rpc.Register(doc)
	rpc.HandleHTTP()

	// Listen for RPC requests on a specific port
	l, err := net.Listen("tcp", ":8079")
	if err != nil {
		log.Fatal(err)
	}

	// Start serving RPC requests
	log.Println("RPC server started, listening on port 8079")
	http.Serve(l, nil)
}
