package main

import "net/http"
import "github.com/justinas/alice"

func (app *application) routes(config Config) http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	/*
		Use the http.NewServeMux() function to initialize a new serveMux, then
		Lets register the home handler for the "/" URL pattern

		Swap the route declarations to use the applications struct methods
	*/
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	/*
		Create a fileServer which serves the static files from ./ui/static directory
	*/
	fileServer := http.FileServer(http.Dir(config.StaticDir))

	/*
		Use mux.Handle to add fileServer as a handle function when the URL
		has /static/ in it (sub tree paths)
	*/
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
