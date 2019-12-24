package main

import (
	"database/sql"
	"flag"
	"github.com/elitenomad/snippetbox/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	_ "github.com/go-sql-driver/mysql"
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
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
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
		DB connection pool
	*/
	dsn := flag.String("dsn", "web:nicetry8@/snippetbox?parseTime=true", "My SQL data source name")

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
		Lets create a OpenDB method which uses the dsn created above to open
		a database connection pool
	 */
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	/*
		Defer a call to db.close()
	 */
	defer db.Close()

	/*
		Initialize the template cache
	 */
	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}


	/*
		Initiailize the application logger
	 */
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	/*
		Initialize a new http.Server struct. We set the Addr and Handler fields so
		that the server uses the same network address and routes as before, and set
		the ErrorLog field so that the server now uses the custom errorLog logger in
		the event of any problems.
	 */
	srv := &http.Server {
		Addr: config.Addr,
		ErrorLog: errorLog,
		Handler: app.routes(*config),
	}

	/*
		Use the http.listenAndServe() function to start a new web server, We pass in two
		paramerters [ Port and mux itself ]
	*/
	infoLog.Printf("Listening on the port %s...", config.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error)  {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
