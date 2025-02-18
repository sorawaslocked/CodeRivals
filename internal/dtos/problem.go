package dtos

import "github.com/sorawaslocked/CodeRivals/internal/entities"

type ProblemCreateRequest struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Difficulty  string            `json:"difficulty"`
	Topics      []*entities.Topic `json:"topics"`
	InputTypes  []string          `json:"input_types"`
	OutputType  string            `json:"output_type"`
	MethodName  string            `json:"method_name"`
	Url         string            `json:"url"`
}
