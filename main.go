package main

import "log"
import "net/http"

/*
	Define a handler function which writes a byte slice containing
	a message as the response body
*/
func home(w http.ResponseWriter, r *http.Request)  {
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

	w.Write([]byte("Hello i am creating a SnippetBox..."))
}

/*
	Define a handler function which will show snippet info 
	Snippet - Show page
*/
func showSnippet(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("Displaying the Snippet info..."))
}

/*
	Define a handler function which creates Snippet
	Snippet - Create form page
*/
func createSnippet(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("Creating a new snippet..."))
}

func main()  {
	/*
		Use the http.NewServeMux() function to initialize a new serveMux, then
		Lets register the home handler for the "/" URL pattern
	*/
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	/*
		Use the http.listenAndServe() function to start a new web server, We pass in two
		paramerters [ Port and mux itself ]
	*/
	log.Println("Listening on the port 4000...")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}