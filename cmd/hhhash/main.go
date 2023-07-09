// This is a program that generates a SHA-256 hash for the headers from an HTTP/HTTPS response.
// The headers are obtained by sending a GET request to a specified URL, and the hash is computed over
// the concatenated header keys.
package main

import (
	"fmt"
	"os"

	"github.com/hrbrmstr/hhhash"
)

// // main function - entry point of the application
func main() {

	// Check if URL is passed as command line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: hhhash <url>")
		os.Exit(1)
	}

	// Make HTTP request, extract headers and generate a hash from headers
	resp := MakeHTTPRequest(os.Args[1])
	headers := ExtractHeaderKeys(resp)
	hhhash := GenerateHHHash(headers)

	fmt.Println(hhhash)

}
