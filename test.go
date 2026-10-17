package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	// "testing"
)

func main() {
	resp, err := http.PostForm("http://localhost:8000/the_good_stuff", url.Values{"name": {"session1"}})
	if err != nil {
		fmt.Println("Connection failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("\n--- DMS RAW RESPONSE ---\n%s\n------------------------\n", string(body))
}