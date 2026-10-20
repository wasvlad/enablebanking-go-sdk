package eb

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type httpClient struct {
	httpClient http.Client
}

func newHttpClient() httpClient {
	return httpClient{
		httpClient: *http.DefaultClient,
	}
}

func (hc *httpClient) R() *request {
	return newRequest(&hc.httpClient)
}

type request struct {
	client  *http.Client
	headers map[string]string
	method  string
	jsonBuf *bytes.Buffer
}

func newRequest(hc *http.Client) *request {
	return &request{
		client:  hc,
		headers: map[string]string{},
	}
}

func (req *request) bearerAuthorization(token string) *request {
	req.headers["Authorization"] = "Bearer " + token
	return req
}

func (req *request) jsonBody(data any) (*request, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(jsonBytes)
	req.jsonBuf = buf
	req.headers["Content-Type"] = "application/json"
	return req, nil
}

func (req *request) post(url string) (*http.Response, error) {
	request, err := http.NewRequest(
		"POST",
		url,
		req.jsonBuf,
	)
	if err != nil {
		return nil, err
	}
	for key, value := range req.headers {
		request.Header.Set(key, value)
	}
	resp, err := req.client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
