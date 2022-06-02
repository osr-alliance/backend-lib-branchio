package branchio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ResponseBodyError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (b *branchio) buildUri(path string, query map[string]interface{}) (*url.URL, error) {
	u, err := url.Parse(b.apiUrl)
	if err != nil {
		return nil, fmt.Errorf("buildUri parse: %v", err)
	}

	u.Path += "/" + path

	if query != nil {
		u.RawQuery = b.buildQueryParams(query)
	}
	return u, nil
}

func (b *branchio) buildQueryParams(query map[string]interface{}) string {
	q := url.Values{}
	for k, v := range query {
		q.Set(k, fmt.Sprintf("%v", v))
	}

	return q.Encode()
}

func (b *branchio) buildHeaders() http.Header {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	return headers
}

func (b *branchio) buildBody(data map[string]interface{}) (io.Reader, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("buildBody json convert: %v", err)
	}
	return bytes.NewBuffer(body), nil
}

func (b *branchio) SendRequest(ctx context.Context, method string, path string, query map[string]interface{}, body map[string]interface{}) (*http.Response, error) {
	req, err := b.buildRequest(ctx, method, path, query, body)
	if err != nil {
		return nil, fmt.Errorf("SendRequest: %v", err)
	}
	return b.branchioClient.Do(req)
}

func (b *branchio) Post(ctx context.Context, path string, body map[string]interface{}) ([]byte, error) {
	resp, err := b.SendRequest(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	bodyResp, err := b.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		err = b.handleResponseError(bodyResp)
		return nil, err
	}

	return bodyResp, nil
}

func (b *branchio) Get(ctx context.Context, path string, query map[string]interface{}) ([]byte, error) {
	resp, err := b.SendRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	bodyResp, err := b.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		err = b.handleResponseError(bodyResp)
		return nil, err
	}

	return bodyResp, nil
}

func (b *branchio) Put(ctx context.Context, path string, body map[string]interface{}, query map[string]interface{}) ([]byte, error) {
	resp, err := b.SendRequest(ctx, http.MethodPut, path, query, body)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	bodyResp, err := b.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		err = b.handleResponseError(bodyResp)
		return nil, err
	}

	return bodyResp, nil
}

func (b *branchio) buildRequest(ctx context.Context, method string, path string, query map[string]interface{}, body map[string]interface{}) (*http.Request, error) {
	//build body
	var bodyReader io.Reader
	var err error
	if method == http.MethodPost || method == http.MethodPut {
		bodyReader, err = b.buildBody(body)
		if err != nil {
			return nil, fmt.Errorf("failed build body: %v", err)
		}
	}
	//build uri
	uri, err := b.buildUri(path, query)
	if err != nil {
		return nil, fmt.Errorf("failed build uri: %v", err)
	}
	//build request
	req, err := http.NewRequestWithContext(ctx, method, uri.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed build request error: %v", err)
	}
	//build headers
	req.Header = b.buildHeaders()
	return req, nil
}

func (b *branchio) readResponseBody(response *http.Response) ([]byte, error) {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response body: %v", err)
	}

	return body, nil
}

func (b *branchio) handleResponseError(bodyResp []byte) error {
	var bodyError ResponseBodyError
	err := json.Unmarshal(bodyResp, &bodyError)
	if err != nil {
		return fmt.Errorf("error parsing response body: %v", err)
	}
	return fmt.Errorf("response error, code: %v: %v", bodyError.Code, bodyError.Message)
}
