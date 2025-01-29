package app

import (
	"github.com/alexedwards/scs/v2"
	"github.com/sorawaslocked/CodeRivals/internal/services"
	"log"
)

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	ProblemService *services.ProblemService
	AuthService    *services.AuthService
	Session        *scs.SessionManager
}
