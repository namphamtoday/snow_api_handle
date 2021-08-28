package model

import "net/url"

type ApiInput struct {
	Url         url.URL           `json:"url" form:"url"`
	Method      string            `json:"method" form:"method"`
	Payload     string            `json:"payload" form:"payload"`
	Headers     map[string]string `json:"headers" form:"headers"`
	QueryParams map[string]string `json:"queryParams" form:"queryParams"`
	FileKey     string            `json:"filekey" form:"filekey"`
	FileName    string            `json:"filename" form:"filename"`
}
