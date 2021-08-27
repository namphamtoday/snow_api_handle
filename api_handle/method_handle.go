package api_handle

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	logs "github.com/sirupsen/logrus"
)

func InitRequest(url url.URL, method string, payload string, headers map[string]string, queryParams map[string]string, filename string) (*http.Request, error) {
	var req *http.Request
	var errInit error
	var body io.Reader

	// set payload
	if payload != "" {
		body = bytes.NewBufferString(payload)
	}

	// if upload file
	if payload == "" && filename != "" {
		body, errInit = os.Open(filename)
		if errInit != nil {
			return nil, errInit
		}
	}

	// not payload and not upload
	if payload == "" && filename == "" {
		body = nil
	}

	req, errInit = http.NewRequest(method, url.String(), body)
	if errInit != nil {
		return nil, errInit
	}

	// set header
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	// set query param
	q := req.URL.Query()
	for key, val := range queryParams {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func Execute(url url.URL, method string, payload string, headers map[string]string, queryParams map[string]string, filename string) ([]byte, error) {
	client := &http.Client{}
	req, err := InitRequest(url, method, payload, headers, queryParams, filename)
	if err != nil {
		logs.Errorf("Cannot create http request %s \n", url.String())
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		logs.Errorf("Error when make request %v\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logs.Errorf("Cannot parse response body %v\n", err)

		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		logs.Errorf("Request returns non OK status: %d\n%s\n", resp.StatusCode, string(body))
		return nil, errors.New("request returns non OK status")
	}

	return body, nil
}

// GET
func Get(url url.URL, payload string, headers map[string]string, queryParams map[string]string, filename string) ([]byte, error) {
	return Execute(url, http.MethodGet, payload, headers, queryParams, filename)
}

// POST
func Post(url url.URL, payload string, headers map[string]string, queryParams map[string]string, filename string) ([]byte, error) {
	return Execute(url, http.MethodPost, payload, headers, queryParams, filename)
}

// PUT
func Put(url url.URL, payload string, headers map[string]string, queryParams map[string]string, filename string) ([]byte, error) {
	return Execute(url, http.MethodPut, payload, headers, queryParams, filename)
}

// DELETE
func Delete(url url.URL, payload string, headers map[string]string, queryParams map[string]string, filename string) ([]byte, error) {
	return Execute(url, http.MethodDelete, payload, headers, queryParams, filename)
}
