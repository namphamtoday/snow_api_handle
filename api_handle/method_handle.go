package api_handle

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/namphamtoday/snow_api_handle/model"
	logs "github.com/sirupsen/logrus"
)

func initRequest(apiInput model.ApiInput) (*http.Request, error) {
	var req *http.Request
	var errInit error
	var body io.Reader

	// set payload
	if apiInput.Payload != "" {
		body = bytes.NewBufferString(apiInput.Payload)
	}

	req, errInit = http.NewRequest(apiInput.Method, apiInput.Url.String(), body)
	if errInit != nil {
		return nil, errInit
	}

	// set header
	for key, val := range apiInput.Headers {
		req.Header.Set(key, val)
	}

	// set query param
	q := req.URL.Query()
	for key, val := range apiInput.QueryParams {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func initRequestUploadFile(apiInput model.ApiInput) (*http.Request, error) {
	fileDir, _ := os.Getwd()
	filePath := path.Join(fileDir, apiInput.FileName)

	file, _ := os.Open(filePath)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(apiInput.FileKey, filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest(apiInput.Method, apiInput.Url.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	// set header
	for key, val := range apiInput.Headers {
		req.Header.Set(key, val)
	}

	// set query param
	q := req.URL.Query()
	for key, val := range apiInput.QueryParams {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()

	return req, nil

}

func Execute(apiInput model.ApiInput) ([]byte, error) {

	var req *http.Request
	var err error

	client := &http.Client{}

	if apiInput.FileKey != "" && apiInput.FileName != "" {
		req, err = initRequestUploadFile(apiInput)
	} else {
		req, err = initRequest(apiInput)
	}

	if err != nil {
		logs.Errorf("Cannot create http request %s \n", apiInput.Url.String())
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
func Get(url url.URL, payload string, headers map[string]string, queryParams map[string]string) ([]byte, error) {

	var apiInput model.ApiInput
	apiInput.Url = url
	apiInput.Payload = payload
	apiInput.Headers = headers
	apiInput.QueryParams = queryParams
	apiInput.Method = http.MethodGet

	return Execute(apiInput)
}

// POST
func Post(url url.URL, payload string, headers map[string]string, queryParams map[string]string) ([]byte, error) {

	var apiInput model.ApiInput
	apiInput.Url = url
	apiInput.Payload = payload
	apiInput.Headers = headers
	apiInput.QueryParams = queryParams
	apiInput.Method = http.MethodPost

	return Execute(apiInput)
}

// PUT
func Put(url url.URL, payload string, headers map[string]string, queryParams map[string]string) ([]byte, error) {
	var apiInput model.ApiInput
	apiInput.Url = url
	apiInput.Payload = payload
	apiInput.Headers = headers
	apiInput.QueryParams = queryParams
	apiInput.Method = http.MethodPut

	return Execute(apiInput)
}

// DELETE
func Delete(url url.URL, payload string, headers map[string]string, queryParams map[string]string) ([]byte, error) {

	var apiInput model.ApiInput
	apiInput.Url = url
	apiInput.Payload = payload
	apiInput.Headers = headers
	apiInput.QueryParams = queryParams
	apiInput.Method = http.MethodDelete

	return Execute(apiInput)
}

// UPLOAD FILE
func Upload(url url.URL, fileKey string, filename string, headers map[string]string, queryParams map[string]string) ([]byte, error) {
	var apiInput model.ApiInput
	apiInput.Url = url
	apiInput.Payload = ""
	apiInput.Headers = headers
	apiInput.QueryParams = queryParams
	apiInput.Method = http.MethodPost
	apiInput.FileKey = fileKey
	apiInput.FileName = filename

	return Execute(apiInput)
}
