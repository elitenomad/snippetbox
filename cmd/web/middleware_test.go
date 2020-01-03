package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T)  {
	/*
		Initialize httptest response recorder
	 */
	record := httptest.NewRecorder()

	/*
		Create a new Dummy request
	 */
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	/*
		Add a dummy next handler
	 */
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})


	/*
		Pass next handler to middleware function
	 */
	secureHeaders(next).ServeHTTP(record, request)

	response := record.Result()

	/*
		Response or input the next middle ware must add all the headers
	 */
	frameOptions := response.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}

	protection := response.Header.Get("X-XSS-Protection")
	if protection != "1; mode=block" {
		t.Errorf("want %q; got %q", "1; mode=block", protection)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected %d; But got %d", http.StatusOK, response.StatusCode)
	}

	/*
		defer
	*/
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("Expected %s; But got %s", "OK", string(body))
	}
}