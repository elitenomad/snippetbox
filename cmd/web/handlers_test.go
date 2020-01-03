package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T)  {
	/*
		Initialize a new Response recorder
	 */
	recorder := httptest.NewRecorder()

	/*
		Initialize a new dummy http.Request
	 */
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	/*
		Call the handler function
	 */
	ping(recorder, request)

	/*
		Call the Result method on Recorder
	 */
	result := recorder.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected %d; But got %d", http.StatusOK, result.StatusCode)
	}

	/*
		defer
	 */
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "Pong" {
		t.Errorf("Expected %s; But got %s", "Pong", string(body))
	}
}
