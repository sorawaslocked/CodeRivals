package main

import (
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
	"log"
)

type application struct {
	errorLog          *log.Logger
	infoLog           *log.Logger
	topicRepository   repositories.TopicRepository
	problemRepository repositories.ProblemRepository
	userRepository    repositories.UserRepository
}
