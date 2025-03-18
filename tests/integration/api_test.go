package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeRequest(method, url string, payload []byte) *http.Request {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestHealthCheck(t *testing.T) {
	resp, err := http.DefaultClient.Do(makeRequest("GET", "http://localhost:8080/health", nil))
	assert.NoError(t, err)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()
	t.Log("health check resp body: ", body)
	assert.Contains(t, body, "Healthy")
}

func TestHandleConfigPUT(t *testing.T) {
	// Successful PUT
	payload := map[string]any{"key": "value"}
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(makeRequest("PUT", "http://localhost:8080/config/123", payloadBytes))
	assert.NoError(t, err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()

	assert.Contains(t, body, "OK: 123")
	resp.Body.Close()

	// Invalid JSON Body
	resp, err = http.DefaultClient.Do(makeRequest("PUT", "http://localhost:8080/config/123", []byte("invalid json")))
	buf = new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	var outdata map[string]any
	json.Unmarshal(buf.Bytes(), &outdata)
	errdata := map[string]any{
		"error": map[string]any{
			"code":    float64(400),
			"message": "could not parse map[string]interface {} from JSON; make sure input has correct types",
		},
	}

	assert.Equal(t, errdata, outdata)
	resp.Body.Close()
}

func TestHandleConfigGET(t *testing.T) {
	payload := map[string]any{"key": "value"}
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(makeRequest("PUT", "http://localhost:8080/config/123", payloadBytes))
	assert.NoError(t, err)

	resp, err = http.DefaultClient.Do(makeRequest("GET", "http://localhost:8080/config/123", nil))
	assert.NoError(t, err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	resp.Body.Close()
	var outdata map[string]any
	json.Unmarshal(buf.Bytes(), &outdata)
	t.Log("BODY: ", outdata)
	errdata := map[string]any{
		"key": "value",
	}
	assert.Equal(t, errdata, outdata)
	resp.Body.Close()
}
