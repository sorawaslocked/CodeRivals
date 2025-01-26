package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdin, "INFO\t", log.Ldate|log.Ltime)

	app := application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	addr := flag.String("addr", ":8080", "Server port")

	flag.Parse()

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.errorLog,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			app.errorLog.Fatal(err)
		}
	}()
}
