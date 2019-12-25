package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			We will write our custom middle ware logic here
			Goal is to add two security headerw
			X-Frame-Options: deny X-XSS-Protection: 1; mode=block
		 */
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		/*
			Create a deferred function which is always run in the event of panic
		 */
		defer func(){
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}