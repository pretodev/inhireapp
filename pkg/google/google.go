package google

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func FetchPage(q string, page int) (string, error) {
	escapedQuery := url.QueryEscape(q)
	googleURL := fmt.Sprintf("https://www.google.com/search?q=%s&start=%d", escapedQuery, page)
	req, err := http.NewRequest("GET", googleURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36")
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request from %s: %v", googleURL, err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read all body: %v", err)
	}
	return string(body), nil
}
