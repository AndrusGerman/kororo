package rest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"kororo/internal/core/ports"
)

type rest struct {
	client *http.Client
}

// Post implements ports.Rest.
func (r *rest) Post(url string, headers map[string]string, body any, out any) error {
	var err error
	var bodyByte []byte
	var resp *http.Response
	var request *http.Request

	if bodyByte, err = json.Marshal(body); err != nil {
		return err
	}

	if request, err = http.NewRequest("POST", url, bytes.NewReader(bodyByte)); err != nil {
		return err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	if resp, err = r.client.Do(request); err != nil {
		return err
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(bodyResp))

	return json.Unmarshal(bodyResp, out)

}

func (r *rest) Stream(url string, headers map[string]string, body any) (<-chan ports.StreamRest, error) {
	var err error
	var bodyByte []byte
	var resp *http.Response
	var scanner *bufio.Scanner
	var streamData = make(chan ports.StreamRest)
	var request *http.Request

	if request, err = http.NewRequest("POST", url, bytes.NewReader(bodyByte)); err != nil {
		return nil, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	if resp, err = r.client.Do(request); err != nil {
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
