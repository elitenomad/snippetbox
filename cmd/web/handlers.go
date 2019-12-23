package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

/*
	Define a handler function which writes a byte slice containing
	a message as the response body
*/
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	/*
		What if "/" must be a strict URL instead of a match all pattern ?
		We will check for request URL path to not be "/" and then use http notFound
		function to send 404 to client. WE RETURN AFTER THAT or else the control will
		still execute and write the byte array to the browser serving the request
	*/
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	/*
		Initialize the handler template files
	 */
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	/*
		Use the template.ParseFiles
	 */

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// We then use the Execute() method on the template set to write the template
	// content as the response body. The last parameter to Execute() represents any
	// dynamic data that we want to pass in, which for now we'll leave as nil.
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
		return
	}
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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

/*
	Define a handler function which creates Snippet
	Snippet - Create form page
*/
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	/*
		Use r.Method to figure out if the request is coming from POST or not.
		http.MethodPost returns a string "POST"
	*/
	if r.Method != http.MethodPost {
		/*
			If its not POST, set the response header to 405 and return
			the control with a message (Subsequent code is not executed
			after return statement.
		*/
		w.Header().Set("Allow", http.MethodPost)

		//w.WriteHeader(405)
		//w.Write([]byte("Method not allowed"))

		/*
			Instead of using writeHeader and write we can use http.Error as well
			to combine both functions
		*/
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

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