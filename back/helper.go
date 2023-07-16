package main

import (
	"bytes"
	"net/http"
)

func SendPostRequest(url string, body []byte, contentType string) (*http.Response, error) {
	client := http.Client{}
	buf := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
