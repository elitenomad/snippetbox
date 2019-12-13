package main

import (
	"flag"
	"log"
	"net/http"
)

type Config struct {
	Addr string
	StaticDir string
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
	fileServer := http.FileServer(http.Dir(config.StaticDir))

	/*
		Use mux.Handle to add fileServer as a handle function when the URL
		has /static/ in it (sub tree paths)
	 */
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	/*
		Use the http.listenAndServe() function to start a new web server, We pass in two
		paramerters [ Port and mux itself ]
	*/
	log.Printf("Listening on the port %s...", config.Addr)
	err := http.ListenAndServe(config.Addr, mux)
	log.Fatal(err)
}
