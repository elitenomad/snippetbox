package main

import (
	"bytes"
	"net/http"
	"testing"
)


func TestPing(t *testing.T)  {
	/*
		Initialize a new Response recorder
	 */
	//recorder := httptest.NewRecorder()

	app := newTestApplication(t)

	/*
		Initialize a new dummy http.Request
	 */
	//request := httptest.NewRequest(http.MethodGet, "/", nil)
	ts := newTestServer(t, app.routes(Config{}))
	defer ts.Close()
	/*
		Call the handler function
	 */
	//ping(recorder, request)
	code, _ , body := ts.get(t, "/ping")

	/*
		Call the Result method on Recorder
	 */
	//result := recorder.Result()

	if code != http.StatusOK {
		t.Errorf("Expected %d; But got %d", http.StatusOK, code)
	}

	if string(body) != "Pong" {
		t.Errorf("Expected %s; But got %s", "Pong", string(body))
	}
}

func TestShowSnippet(t *testing.T) {
	/*
		Create a new instance of our application struct which uses the mocked // dependencies.
	 */
	app := newTestApplication(t)
	/*
		Establish a new test server for running end-to-end tests.
	 */
	ts := newTestServer(t, app.routes(Config{}))
	defer ts.Close()

	/*
		Set up some table-driven tests to check the responses sent by our // application for different URLs.
	 */
	tests := []struct {
		name string
		urlPath string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests { t.Run(tt.name, func(t *testing.T) {
		code, _, body := ts.get(t, tt.urlPath)

		if code != tt.wantCode {
			t.Errorf("want %d; got %d", tt.wantCode, code)
		}

		if !bytes.Contains(body, tt.wantBody) {
			t.Errorf("want body to contain %q", tt.wantBody)
		} })
	}
}