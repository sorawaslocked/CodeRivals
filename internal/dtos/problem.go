package dtos

import "github.com/sorawaslocked/CodeRivals/internal/entities"

type ProblemCreateRequest struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Difficulty  string            `json:"difficulty"`
	Topics      []*entities.Topic `json:"topics"`
}
