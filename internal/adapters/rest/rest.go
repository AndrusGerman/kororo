package rest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"

	"kororo/internal/core/ports"
)

type rest struct {
	client *http.Client
}

// Post implements ports.Rest.
func (r *rest) Post(url string, body any, out any) error {
	var err error
	var bodyByte []byte
	var resp *http.Response

	if bodyByte, err = json.Marshal(body); err != nil {
		return err
	}

	if resp, err = r.client.Post(url, "application/json", bytes.NewReader(bodyByte)); err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func (r *rest) Stream(url string, body any) (<-chan ports.StreamRest, error) {
	var err error
	var bodyByte []byte
	var resp *http.Response
	var scanner *bufio.Scanner
	var streamData = make(chan ports.StreamRest)

	if bodyByte, err = json.Marshal(body); err != nil {
		return nil, err
	}

	if resp, err = r.client.Post(url, "application/json", bytes.NewReader(bodyByte)); err != nil {
		return nil, err
	}

	scanner = bufio.NewScanner(resp.Body)
	go func() {
		for scanner.Scan() {
			streamData <- NewStreamRest(scanner.Bytes())
		}
		close(streamData)
	}()

	return streamData, nil
}

func New() ports.RestAdapter {
	return &rest{
		client: &http.Client{},
	}
}
