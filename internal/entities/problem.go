package entities

import "time"

type Problem struct {
	ID          int              `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Difficulty  string           `json:"difficulty"`
	Url         string           `json:"url"`
	Topics      []*Topic         `json:"topics"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Examples    []ProblemExample `json:"examples"`
}
