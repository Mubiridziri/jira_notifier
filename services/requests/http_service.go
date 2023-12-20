package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func MakeFormData(url string, form map[string]string, fileFieldName string, files map[string]string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, filepath := range files {
		file, _ := os.Open(filepath)
		part, _ := writer.CreateFormFile("photo", filepath)
		io.Copy(part, file)
		defer file.Close()
	}

	for field, value := range form {
		part, _ := writer.CreateFormField(field)
		io.Copy(part, strings.NewReader(value))
	}

	writer.Close()

	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Check the response
	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respBytes := buf.String()
		return errors.New(respBytes)
	}

	return nil

}

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
