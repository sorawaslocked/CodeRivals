package app

import (
	"github.com/sorawaslocked/CodeRivals/internal/services"
	"log"
)

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	ProblemService *services.ProblemService
	AuthService    *services.AuthService
}
