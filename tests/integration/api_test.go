package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ACED-IDP/gecko/gecko/config"
	"github.com/ACED-IDP/gecko/tests/fixtures"
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
	var configs []config.ConfigItem
	err := json.Unmarshal([]byte(fixtures.TestConfig), &configs)
	t.Log("CONFIGS: ", configs)
	assert.NoError(t, err)
	marshalledJSON, err := json.Marshal(configs)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(makeRequest("PUT", "http://localhost:8080/config/123", marshalledJSON))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	assert.NoError(t, err)

	var outData map[string]any
	err = json.Unmarshal(buf.Bytes(), &outData)
	assert.NoError(t, err)
	t.Log("RESP: ", outData)

	expected200Response := map[string]any{
		"code": float64(200), "message": "ACCEPTED: 123",
	}
	assert.Equal(t, expected200Response, outData)
}

func TestHandleConfigPUTInvalidJson(t *testing.T) {
	resp, err := http.DefaultClient.Do(makeRequest("PUT", "http://localhost:8080/config/123", []byte("invalid json")))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	assert.NoError(t, err)

	var errData map[string]any
	err = json.Unmarshal(buf.Bytes(), &errData)
	assert.NoError(t, err)

	expectedErrorResponse := map[string]any{
		"error": map[string]any{
			"code":    float64(400),
			"message": "could not parse []config.ConfigItem from JSON; make sure input has correct types",
		},
	}
	assert.Equal(t, expectedErrorResponse, errData)
}

func TestHandleConfigGET(t *testing.T) {
	var configs []config.ConfigItem
	err := json.Unmarshal([]byte(fixtures.TestConfig), &configs)

	payloadBytes, err := json.Marshal(configs)
	assert.NoError(t, err)

	_, err = http.DefaultClient.Do(makeRequest("PUT", "http://localhost:8080/config/123", payloadBytes))
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(makeRequest("GET", "http://localhost:8080/config/123", nil))
	assert.NoError(t, err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	resp.Body.Close()
	var outdata map[string]any
	json.Unmarshal(buf.Bytes(), &outdata)

	var Resconfigs []config.ConfigItem
	data, _ := json.Marshal(outdata["content"])
	json.Unmarshal(data, &Resconfigs)

	assert.Equal(t, configs, Resconfigs)
	resp.Body.Close()
}

func TestHandle404ConfigGet(t *testing.T) {
	resp, err := http.DefaultClient.Do(makeRequest("GET", "http://localhost:8080/config/404config", nil))
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 404)
}

func TestHandle404ConfigDelete(t *testing.T) {
	resp, err := http.DefaultClient.Do(makeRequest("DELETE", "http://localhost:8080/config/404config", nil))
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 404)
}

func TestHandleConfigDeleteOK(t *testing.T) {
	var configs []config.ConfigItem
	err := json.Unmarshal([]byte(fixtures.TestConfig), &configs)
	payloadBytes, err := json.Marshal(configs)
	assert.NoError(t, err)
	_, err = http.DefaultClient.Do(makeRequest("PUT", "http://localhost:8080/config/testdelete", payloadBytes))
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(makeRequest("DELETE", "http://localhost:8080/config/testdelete", nil))
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)

	resp, err = http.DefaultClient.Do(makeRequest("GET", "http://localhost:8080/config/testdelete", nil))
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 404)
}
