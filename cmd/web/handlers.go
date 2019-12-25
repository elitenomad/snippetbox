package main

import (
	"errors"
	"fmt"
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
		Pass the dummy data
	 */
	title := "Pranava S Balugari"
	content := "He is hard working bloke who is constantly planning to improve himself"
	expires := "7"

	/*
		execute the snippets insert with the data collected
	 */
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	/*
		Redirect the User to the relavant Snippet page
	 */
	http.Redirect(w,r,fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}


// Add a new createSnippetForm handler, which for now returns a placeholder response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet..."))
}