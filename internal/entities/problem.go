package entities

import "time"

type Problem struct {
	ID          int               `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Difficulty  string            `json:"difficulty"`
	Url         string            `json:"url"`
	Topics      []*Topic          `json:"topics"`
	InputTypes  []string          `json:"input_types"`
	OutputType  string            `json:"output_type"`
	MethodName  string            `json:"method_name"`
	Examples    []*ProblemExample `json:"examples"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
