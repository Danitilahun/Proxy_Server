package main

import (
	"log"
	"net/rpc"
)

type GetURLRequest struct {
	URL string `json:"url"`
}

type GetURLResponse struct {
	Content string `json:"content"`
}

type GetStatsResponse struct {
	Stats map[string]int64 `json:"stats"`
}

func main() {
	// Connect to the RPC server
	client, err := rpc.DialHTTP("tcp", "localhost:8079")
	if err != nil {
		log.Fatal(err)
	}

	// Create the request object
	request := GetURLRequest{
		URL: "https://golang.org/doc",
	}

	// Call the remote method to get the URL content
	var response GetURLResponse
	err = client.Call("Documentation.GetURL", request, &response)
	if err != nil {
		log.Fatal(err)
	}

	// Print the fetched content
	log.Println("Fetched content:")
	log.Println(response.Content)

	// Call the remote method to get the request statistics
	var statsResponse GetStatsResponse
	err = client.Call("Documentation.GetStats", struct{}{}, &statsResponse)
	if err != nil {
		log.Fatal(err)
	}

	// Print the request statistics
	log.Println("Request statistics:")
	for url, bytes := range statsResponse.Stats {
		log.Printf("URL: %s, Bytes: %d\n", url, bytes)
	}
}
