package entities

import "time"

type ProblemSolutionComment struct {
	ID         int
	UserID     int
	Username   string
	SolutionID int
	CommentID  *int
	TextValue  string
	CreatedAt  time.Time
}
