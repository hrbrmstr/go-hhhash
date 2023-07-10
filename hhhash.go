// This is a program that generates a SHA-256 hash for the headers from an HTTP/HTTPS response.
// The headers are obtained by sending a GET request to a specified URL, and the hash is computed over
// the concatenated header keys.
package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

// MakeHTTPRequest function takes a raw URL string as input,
// parses the URL and establishes a TCP connection based on the scheme (HTTP or HTTPS).
// It then sends a GET request and returns the response as a byte array.
func MakeHTTPRequest(rawURL string) ([]byte, error) {

	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	// Default port based on the URL scheme
	port := parsedURL.Port()

	if parsedURL.Scheme == "http" {
		if port == "" {
			port = "80"
		}
		conn, _ := net.Dial("tcp", parsedURL.Host+":"+port)
		return handleConnection(conn, err, parsedURL.Host, parsedURL.Path), nil
	} else if parsedURL.Scheme == "https" {
		if port == "" {
			port = "443"
		}
		conn, _ := net.Dial("tcp", parsedURL.Host+":"+port)
		tlsConfig := &tls.Config{
			ServerName:         parsedURL.Host,
			InsecureSkipVerify: true,
		}
		tlsConn := tls.Client(conn, tlsConfig)
		err = tlsConn.Handshake()
		return handleConnection(tlsConn, err, parsedURL.Host, parsedURL.Path), nil
	} else {
		return nil, errors.New("Invalid URL")
	}
}

// handleConnection function writes the GET request to the established connection
// and reads the response. The response is then returned as a byte array.
func handleConnection(conn net.Conn, err error, host string, path string) []byte {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return []byte{}
	}
	defer conn.Close()

	request := "GET / HTTP/1.1\r\nHost: " + host + "\r\n\r\n"
	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return []byte{}
	}

	reader := bufio.NewReader(conn)
	buf := make([]byte, 8192)
	bytesRead, _ := reader.Read(buf)

	return buf[:bytesRead]

}

// ExtractHeaderKeys function takes in the response headers as a byte array,
// and returns a slice of header keys.
func ExtractHeaderKeys(resp []byte) []string {

	lines := strings.Split(string(resp), "\n")

	firstEmptyIndex := -1
	for i, str := range lines {
		if str == "" {
			firstEmptyIndex = i
			break
		} else {
			ind := strings.Index(str, ":")
			if ind != -1 {
				lines[i] = str[:ind]
			}
		}
	}

	return lines[1:(firstEmptyIndex - 1)]

}

// GenerateHHHash function takes in a slice of header keys,
// concatenates the elements together with a :, and computes a SHA256 hash of it.
// The hash is then returned as a string.
func GenerateHHHash(headers []string) string {
	data := strings.Join(headers, ":")
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("hhh:1:%x", hash)
}

// function checks given string is a valid URL or not
func isValidURL(url string) bool {
	// check if the string starts with http:// or https://
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}

	return true
}

func main() {

	// Check if URL is passed as command line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: hhhash <url>")
		os.Exit(1)
	}

	if !isValidURL(os.Args[1]) {
		fmt.Println("Given URL is not valid\nCorrect format: http(s)://<domain>/")
		os.Exit(1)
	}

	// Make HTTP request, extract headers and generate a hash from headers
	resp, err := MakeHTTPRequest(os.Args[1])
	if err != nil {
		fmt.Printf("error occured while making HTTP request: %v\n", err.Error())
		os.Exit(1)
	}
	headers := ExtractHeaderKeys(resp)
	hhhash := GenerateHHHash(headers)

	fmt.Println(hhhash)

}
