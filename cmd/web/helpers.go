package main

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
	"runtime/debug"
	"time"
)

/*
	The serverError helper writes an error message and stack trace to the errorLog,
	then sends a generic 500 Internal Server Error response to the user.
 */
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

/*
	The clientError helper sends a specific status code and corresponding description
	to the user. We'll use this later in the book to send responses like 400 "Bad
	Request" when there's a problem with the request that the user sent.
 */
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

/*
	For consistency, we'll also implement a notFound helper. This is simply a
	convenience wrapper around clientError which sends a 404 Not Found response to
 	the user.
 */
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) injectDefaultData(data *templateData, r *http.Request) *templateData {
	if data == nil {
		data = &templateData{}
	}

	data.CurrentYear = time.Now().Year()
	/*
		Show the flash if exists
	*/
	data.Flash = app.session.PopString(r, "flash")
	data.IsAuthenticated = app.isAuthenticated(r)
	data.CSRFToken = nosurf.Token(r)

	return data
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the
	// provided name, call the serverError helper method that we made earlier.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	/*
		Initialize the buffer
	*/
	buf := new(bytes.Buffer)

	// Execute the template set, passing in any dynamic data.
	err := ts.Execute(buf, app.injectDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
	}

	buf.WriteTo(w)
}

/*
	Return true if the current request is from authenticated user, otherwise return false.
 */
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}