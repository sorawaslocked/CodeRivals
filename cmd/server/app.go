package main

import (
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
	"log"
)

type application struct {
	errorLog          *log.Logger
	infoLog           *log.Logger
	problemRepository repositories.ProblemRepository
}
