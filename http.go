package eb

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const contentType = "Content-Type"
const applicationJson = "application/json"

var notJsonBodyError = errors.New("Content is not json")

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

type response struct {
	response *http.Response
	body     *[]byte
}

func newResponse(res *http.Response, body *[]byte) *response {
	return &response{
		response: res,
		body:     body,
	}
}

func (res *response) isSuccess() bool {
	return res.response.Status[0] == '2'
}

func (res *response) isError() bool {
	return res.response.Status[0] == '4'
}

func (res *response) isJson() bool {
	return res.response.Header.Get(contentType) == applicationJson
}

func (res *response) parseJson(into any) error {
	if !res.isJson() {
		return notJsonBodyError
	}
	reader := bytes.NewReader(*res.body)
	return json.NewDecoder(reader).Decode(into)
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
	req.headers[contentType] = applicationJson
	return req, nil
}

func (req *request) post(url string) (*response, error) {
	if req.jsonBuf == nil {
		req.jsonBuf = new(bytes.Buffer)
	}
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
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return newResponse(resp, &body), nil
}
