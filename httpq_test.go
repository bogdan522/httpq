package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublish(t *testing.T) {
	h := &HTTPQ{}

	msg := []byte("Hello")
	topic := "test-topic"

	req, err := http.NewRequest("POST", "/"+topic, bytes.NewBuffer(msg))
	if err != nil {
		t.Errorf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	h.Publish().ServeHTTP(recorder, req)

	if recorder.Code != 201 {
		t.Errorf("Expected status code 201, got %d", recorder.Code)
	}

	if recorder.Body.String() != "" {
		t.Errorf("Expected empty string in response body, got '%s'", recorder.Body.String())
	}
}

func TestPublishFail(t *testing.T) {
	h := &HTTPQ{}
	topic := "test-topic"

	req, err := http.NewRequest("POST", "/"+topic, bytes.NewBuffer(nil))
	if err != nil {
		t.Errorf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	h.Publish().ServeHTTP(recorder, req)

	if recorder.Code != 400 {
		t.Errorf("Expected status code 400, got %d", recorder.Code)
	}

	if h.PubFails != 1 {
		t.Errorf("Expected PubFails to be 1, got %d", h.PubFails)
	}
}

func TestConsume(t *testing.T) {
	h := &HTTPQ{queue: map[string][]string{"/test-topic": {"Hello"}}}

	req, err := http.NewRequest("GET", "/test-topic", nil)

	if err != nil {
		t.Errorf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	h.Consume().ServeHTTP(recorder, req)

	if recorder.Code != 200 {
		t.Errorf("Expected status code 200, got %d", recorder.Code)
	}

	if h.SubFails > 0 {
		t.Errorf("Expected SubFails to be 0, got %d", h.SubFails)
	}

	if recorder.Body.String() == "" {
		t.Error("Expected non-empty response body")
	}
}
func TestConsumeFail(t *testing.T) {
	h := &HTTPQ{}
	req, err := http.NewRequest("GET", "/test-topic", nil)

	if err != nil {
		t.Errorf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	h.Consume().ServeHTTP(recorder, req)

	if recorder.Code != 400 {
		t.Errorf("Expected status code 400, got %d", recorder.Code)
	}

	if h.SubFails != 1 {
		t.Errorf("Expected SubFails to be 1, got %d", h.SubFails)
	}

	if recorder.Body.String() != "" {
		t.Errorf("Expected empty response body, got %s", recorder.Body.String())
	}
}

func TestPublishAndConsume(t *testing.T) {
	h := &HTTPQ{}

	msg := []byte("Hello")
	topic := "test-topic"

	postReq, err := http.NewRequest("POST", "/"+topic, bytes.NewBuffer(msg))

	if err != nil {
		t.Errorf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	h.Publish().ServeHTTP(recorder, postReq)

	getReq, err := http.NewRequest("GET", "/"+topic, nil)

	if err != nil {
		t.Errorf("Failed to create request: %v", err)
	}

	h.Consume().ServeHTTP(recorder, getReq)

	if recorder.Code != 201 {
		t.Errorf("Expected status code 201, got %d", recorder.Code)
	}

	if !(h.PubFails == 0) {
		t.Errorf("Expected PubFails to be 0, got %d", h.PubFails)
	}

	if !(h.SubFails == 0) {
		t.Errorf("Expected SubFails to be 0, got %d", h.PubFails)
	}

	if h.RxBytes != h.TxBytes {
		t.Errorf("Expected to consume %d. Actually consumed %d", h.TxBytes, h.RxBytes)
	}

}
