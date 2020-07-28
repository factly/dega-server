package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Request - test request
func Request(t *testing.T, ts *httptest.Server, method, path string, body io.Reader, header map[string]string) (*http.Response, map[string]interface{}, int) {

	// err := godotenv.Load()
	// if err != nil {
	// 	t.Errorf("%v", err)
	// }

	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, nil, http.StatusServiceUnavailable
	}

	req.Header = map[string][]string{
		"X-User":  {header["user"]},
		"X-Space": {header["space"]},
		"Cookie":  {header["cookie"]},
	}

	// req.Header = header

	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, nil, http.StatusServiceUnavailable
	}

	var respBody map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&respBody)
	defer resp.Body.Close()

	return resp, respBody, resp.StatusCode
}
