package app

import (
	"github.com/alexedwards/scs/v2"
	"github.com/sorawaslocked/CodeRivals/internal/services"
	"html/template"
	"log"
)

type Application struct {
	ErrorLog                      *log.Logger
	InfoLog                       *log.Logger
	ProblemService                *services.ProblemService
	AuthService                   *services.AuthService
	SubmissionService             *services.SubmissionService
	CodeExecutionService          *services.CodeExecutionService
	TopicService                  *services.TopicService
	LeaderBoardService            *services.LeaderboardService
	Session                       *scs.SessionManager
	templateCache                 map[string]*template.Template
	CommentService                *services.CommentService
	ProblemSolutionCommentService *services.SolutionCommentService
}
