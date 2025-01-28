package app

import (
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
	"log"
)

type Application struct {
	ErrorLog          *log.Logger
	InfoLog           *log.Logger
	TopicRepository   repositories.TopicRepository
	ProblemRepository repositories.ProblemRepository
	UserRepository    repositories.UserRepository
}
