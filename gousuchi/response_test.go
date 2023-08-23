package gousuchi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResponse(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	resp := NewResponse(req, http.StatusAccepted, ContentTypeApplicationJSON, []byte("{}"))

	assert.NotNil(t, resp)

	writter := httptest.NewRecorder()

	errResp := resp.Write(writter)

	assert.Nil(t, errResp)
	assert.Equal(t, http.StatusAccepted, writter.Result().StatusCode)
	assert.Equal(t, "application/json", writter.Header().Get("Content-Type"))
	assert.Equal(t, []byte("{}"), writter.Body.Bytes())
}

func TestJSON(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	resp := JSON(req, map[string]interface{}{"test": "world"})

	assert.NotNil(t, resp)

	writter := httptest.NewRecorder()

	errResp := resp.Write(writter)

	assert.Nil(t, errResp)
	assert.Equal(t, http.StatusOK, writter.Result().StatusCode)
	assert.Equal(t, "application/json", writter.Header().Get("Content-Type"))
	assert.Equal(t, []byte("{\"test\":\"world\"}"), writter.Body.Bytes())
}

func TestWithHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	headers := http.Header{}

	headers.Add("X-Test", "value1")
	headers.Add("X-Test", "value2")

	resp := NewResponse(req, http.StatusCreated, ContentTypeApplicationJSON, []byte("{}")).
		WithHeaders(headers)

	assert.NotNil(t, resp)

	writter := httptest.NewRecorder()

	errResp := resp.Write(writter)

	assert.Nil(t, errResp)
	assert.Equal(t, http.StatusCreated, writter.Result().StatusCode)
	assert.Equal(t, "application/json", writter.Header().Get("Content-Type"))
	assert.Equal(t, []string{"value1", "value2"}, writter.Header().Values("X-Test"))
	assert.Equal(t, []byte("{}"), writter.Body.Bytes())
}
