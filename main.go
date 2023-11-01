package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	// Check if a URL argument is provided
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <URL>")
		return
	}

	// Get the URL from the command line argument
	url := os.Args[1]

	// Send an HTTP GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer response.Body.Close()

	// Read the response body (HTML content)
	htmlContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	// Parse the HTML content
	htmlReader := strings.NewReader(string(htmlContent))
	tokenizer := html.NewTokenizer(htmlReader)

	// Create a slice to store extracted links
	links := []string{}

	// Create a map to track duplicates
	seenLinks := make(map[string]struct{})

	for tokenType := tokenizer.Next(); tokenType != html.ErrorToken; {
		switch tokenType {
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				// Check for anchor tags
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						// Check for duplicates
						if _, seen := seenLinks[attr.Val]; !seen {
							// Found a new link, add it to the list
							links = append(links, attr.Val)
							seenLinks[attr.Val] = struct{}{}
						}
					}
				}
			}
		}

		// Get the next token
		tokenType = tokenizer.Next()
	}

	// Print the extracted links
	for _, link := range links {
		// fmt.Printf("Link %d: %s\n", i+1, link)
		fmt.Printf("%s\n", link)
	}
}

