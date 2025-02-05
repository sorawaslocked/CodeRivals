package main

import (
	"database/sql"
	"errors"
	"flag"
	"github.com/alexedwards/scs/v2"
	_ "github.com/lib/pq"
	"github.com/sorawaslocked/CodeRivals/internal/app"
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
	"github.com/sorawaslocked/CodeRivals/internal/services"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	db, err := openDB("user=postgres password=1234 dbname=coderivals sslmode=disable host=localhost")

	addr := flag.String("addr", ":8080", "HTTP address")

	if err != nil {
		errorLog.Print("Failed to connect to database")
		errorLog.Fatal(err)
	}

	topicRepository := repositories.NewPGTopicRepository(db)
	problemRepository := repositories.NewPGProblemRepository(db, topicRepository)
	userRepository := repositories.NewPGUserRepository(db)
	leaderboardService := services.NewLeaderboardService(userRepository)

	problemService := services.NewProblemService(problemRepository)
	authService := services.NewAuthService(userRepository)
	codeExecutionService := services.NewCodeExecutionService()
	topicService := services.NewTopicService(topicRepository)

	session := scs.New()
	session.Lifetime = 24 * time.Hour

	app := app.Application{
		ErrorLog:             errorLog,
		InfoLog:              infoLog,
		ProblemService:       problemService,
		AuthService:          authService,
		CodeExecutionService: codeExecutionService,
		TopicService:         topicService,
		Session:              session,
		LeaderBoardService:   leaderboardService,
	}

	// Initialize templates
	err = app.InitTemplates()
	if err != nil {
		errorLog.Fatal(err)
	}

	flag.Parse()

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.ErrorLog,
		Handler:  app.Routes(),
	}

	infoLog.Printf("Server started on %s", *addr)
	err = srv.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.ErrorLog.Fatal(err)
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
