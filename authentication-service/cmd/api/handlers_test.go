package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RoundTripFunc func(request *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func Test_Authenticate(t *testing.T) {
	jsonToReturn := `
{
	"error": "false",
	"message": "logged",
}
`
	client := NewTestClient(func(request *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewBufferString(jsonToReturn)),
		}
	})

	testApp.Client = client

	postBody := map[string]any{
		"email":    "admin@admin.com",
		"password": "password",
	}
	body, _ := json.Marshal(postBody)
	req, _ := http.NewRequest("POST", "/authenticate", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.Authenticate)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status 202, but got %d", rr.Code)
	}
}
