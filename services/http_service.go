package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func MakeGetJsonRequest(url string, headers map[string]string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for header, value := range headers {
		req.Header.Set(header, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(resp.StatusCode))
	}

	jsonMap, err := decodeJsonResponse(resp)

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

func MakeGetRequest(url string, headers map[string]string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for header, value := range headers {
		req.Header.Set(header, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(resp.StatusCode))
	}

	jsonMap, err := decodeJsonResponse(resp)

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

func MakePostJsonRequest(url string, body map[string]interface{}, headers map[string]string) (map[string]interface{}, error) {
	jsonStr, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	for header, value := range headers {
		req.Header.Set(header, value)
	}

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(resp.StatusCode))
	}

	jsonMap, err := decodeJsonResponse(resp)

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

func decodeJsonResponse(resp *http.Response) (map[string]interface{}, error) {
	b, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var jsonMap map[string]interface{}
	err = json.Unmarshal(b, &jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}
