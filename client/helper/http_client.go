package helper

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	errUnexpectedResponse = "unexpected response: %s"
)

type HTTPClient struct{}

var (
	HttpClient = HTTPClient{}
)

var backoffSchedule = []time.Duration{
	10 * time.Second,
	15 * time.Second,
	20 * time.Second,
	25 * time.Second,
	30 * time.Second,
	40 * time.Second,
	50 * time.Second,
	1 * time.Minute,
	90 * time.Second,
	2 * time.Minute,
}

func (c HTTPClient) PostRequest(uri string, buf []byte, contentType string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", uri, bytes.NewReader(buf))
	req.Header.Set("Content-Type", contentType)
	time.Sleep(500 * time.Millisecond)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	c.info(fmt.Sprintf("POST %s -> %d", req.URL, resp.StatusCode))

	if resp.StatusCode != 200 {
		respErr := fmt.Errorf(errUnexpectedResponse, resp.Status)
		fmt.Sprintf("request failed: %v", respErr)
		return nil, respErr
	}
	return resp, nil
}

func (c HTTPClient) PostRequestWithRetries (uri string, buf []byte, contentType string) (*http.Response, error) {
	var resp *http.Response
	var err error
	for _, backoff := range backoffSchedule {
		resp, err = c.PostRequest(uri, buf, contentType)
		if err == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "Request error: %+v\n", err)
		fmt.Fprintf(os.Stderr, "Retrying in %v\n", backoff)
		time.Sleep(backoff)
	}

	return resp, nil
}

func (c HTTPClient) info(msg string) {
	log.Printf("[client] %s\n", msg)
}