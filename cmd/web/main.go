package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Addr string
	StaticDir string
}

/*
	Define an application struct which holds application wide
	dependencies.
 */
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	/*
		Define a new flag with a name addr which take a string of format ":{PORT_NUMBER}"
		and add some text to help explaning what the command-line flag does
	 */
	config := new(Config)
	flag.StringVar(&config.Addr, "addr", ":4000", "Port number on which SnippetBox webserver runs")
	flag.StringVar(&config.StaticDir,  "static-dir", "./ui/static", "Static files directory")

	/*
		We need to use flag.Parse to parse the command Line flag
	 */
	flag.Parse()

	/*
		Logging
	 */
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	/*
		Initiailize the application logger
	 */
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

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


	/*
		Initialize a new http.Server struct. We set the Addr and Handler fields so
		that the server uses the same network address and routes as before, and set
		the ErrorLog field so that the server now uses the custom errorLog logger in
		the event of any problems.
	 */
	srv := &http.Server {
		Addr: config.Addr,
		ErrorLog: errorLog,
		Handler: mux,
	}

	/*
		Use the http.listenAndServe() function to start a new web server, We pass in two
		paramerters [ Port and mux itself ]
	*/
	infoLog.Printf("Listening on the port %s...", config.Addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
