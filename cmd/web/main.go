package main

import (
	"log"
	"net/http"
)

func main() {
	/*
		Use the http.NewServeMux() function to initialize a new serveMux, then
		Lets register the home handler for the "/" URL pattern
	*/
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	/*
		Create a fileServer which serves the static files from ./ui/static directory
	 */
	fileServer := http.FileServer(http.Dir("./ui/static"))

	/*
		Use mux.Handle to add fileServer as a handle function when the URL
		has /static/ in it (sub tree paths)
	 */
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	/*
		Use the http.listenAndServe() function to start a new web server, We pass in two
		paramerters [ Port and mux itself ]
	*/
	log.Println("Listening on the port 4000...")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
