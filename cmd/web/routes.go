package main

import (
	"github.com/bmizerany/pat"
	"net/http"
)
import "github.com/justinas/alice"

func (app *application) routes(config Config) http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	/*
		Use the http.NewServeMux() function to initialize a new serveMux, then
		Lets register the home handler for the "/" URL pattern

		Swap the route declarations to use the applications struct methods
	*/
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.home)))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(http.HandlerFunc(app.createSnippetForm)))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(http.HandlerFunc(app.createSnippet)))
	mux.Get("/snippet/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(http.HandlerFunc(app.showSnippet)))
	/*
		User Sign up,in and out routes
	 */
	mux.Get("/users/signup", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.signupUserForm)))
	mux.Post("/users/signup", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.signupUser)))
	mux.Get("/users/login", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.loginUserForm)))
	mux.Post("/users/login", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.loginUser)))
	mux.Post("/users/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(http.HandlerFunc(app.logoutUser)))

	mux.Get("/ping", http.HandlerFunc(ping))
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
