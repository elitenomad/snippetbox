package main

import (
	"github.com/bmizerany/pat"
	"net/http"
)
import "github.com/justinas/alice"

func (app *application) routes(config Config) http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	/*
		Use the http.NewServeMux() function to initialize a new serveMux, then
		Lets register the home handler for the "/" URL pattern

		Swap the route declarations to use the applications struct methods
	*/
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	/*
		Create a fileServer which serves the static files from ./ui/static directory
	*/
	fileServer := http.FileServer(http.Dir(config.StaticDir))

	/*
		Use mux.Handle to add fileServer as a handle function when the URL
		has /static/ in it (sub tree paths)
	*/
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
