package main

import (
	"errors"
	"fmt"
	"github.com/elitenomad/snippetbox/pkg/forms"
	"github.com/elitenomad/snippetbox/pkg/models"
	"net/http"
	"strconv"
)

/*
	Define a handler function which writes a byte slice containing
	a message as the response body
*/
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	/*
		Collect all snippets limited by 10 in Latest method
	*/
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: snippets,
	})
}

/*
	Define a handler function which will show snippet info
	Snippet - Show page
*/
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	/*
		Extrac the id value from the query string. If there is a error
		or id value is < 1 we return a NotFound error on http module
	*/
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	/*
		Fetch the snippet by id
	 */
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		}else{
			app.serverError(w, err)
		}
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: snippet,
	})
}

/*
	Define a handler function which creates Snippet
	Snippet - Create form page
*/
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	/*
		ParseForm
	 */
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	/*
		Fetch the form info
	 */
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	/*
		Initialize errors
	 */

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	/*
		execute the snippets insert with the data collected
	 */
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	/*
		Add a flash to sessions string key
	 */
	app.session.Put(r, "flash", "Snippet successfully created!")

	/*
		Redirect the User to the relavant Snippet page
	 */
	http.Redirect(w,r,fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}


// Add a new createSnippetForm handler, which for now returns a placeholder response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{Form: forms.New(nil)})
}

/*
	Add a Signup user form Handler
 */
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{Form: forms.New(nil)})
}

/*
	Add a Signup form submission handler function
 */
func (app *application) signupUser(w http.ResponseWriter, r *http.Request)  {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 200)
	form.MaxLength("email",200)
	form.MatchesPattern("email", forms.RxEmail)
	form.MinLength("password", 10)

	if !form.Valid() {	
		app.render(w, r,"signup.page.tmpl", &templateData{
			Form: form,
		})

		return
	}

}

/*
	Add a Login user form
 */
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request)  {
	app.render(w, r, "login.page.tmpl", &templateData{Form: forms.New(nil)})
}

/*
	Authenticate the user
 */
func (app *application) loginUser(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintln(w, "Authenticating the user")
}

/*
	Destroying the session
*/
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintln(w, "Destroying the user session")
}