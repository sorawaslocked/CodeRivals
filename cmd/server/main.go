package main

import (
	"database/sql"
	"errors"
	"flag"
	_ "github.com/lib/pq"
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
	"log"
	"net/http"
	"os"
)

func main() {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	db, err := openDB("user=postgres password=1234 dbname=coderivals sslmode=disable host=localhost")

	if err != nil {
		errorLog.Print("Failed to connect to database")
		errorLog.Fatal(err)
	}

	topicRepository := repositories.NewPGTopicRepository(db)
	problemRepository := repositories.NewPGProblemRepository(db, topicRepository)

	app := application{
		errorLog:          errorLog,
		infoLog:           infoLog,
		topicRepository:   topicRepository,
		problemRepository: problemRepository,
	}

	addr := flag.String("addr", ":8080", "Server port")

	flag.Parse()

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.errorLog,
	}

	infoLog.Printf("Server started on %s", *addr)
	err = srv.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.errorLog.Fatal(err)
	}
}

func openDB(connectionStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
