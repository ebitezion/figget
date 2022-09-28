package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

//version
const version = "1.0.0"
const cssVersion = "1"

//config
type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

//application-wide specific config resource
type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
}

//serve: the web server
func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
	app.infoLog.Printf("Starting Server on PORT: %d and %s environment", app.config.port, app.config.env)
	return srv.ListenAndServe()
}
func main() {
	var cfg config

	//commandline flag
	flag.IntVar(&cfg.port, "port", 4014, "Application Server Port address")
	flag.StringVar(&cfg.env, "env", "development", "Application environment{development|production}")
	flag.StringVar(&cfg.api, "api", "http://localhost:4014", "Application API")
	flag.Parse()

	//read secret key from env var
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	//set-up Logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//set-up template cache
	tc := make(map[string]*template.Template)

	//init application
	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
	}
	err := app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}

}
