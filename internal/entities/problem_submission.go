package entities

import "time"

type ProblemSubmission struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ProblemID   int       `json:"problem_id"`
	Code        string    `json:"code"`
	Status      string    `json:"status"`
	Runtime     uint32    `json:"runtime_ms"`
	Memory      uint32    `json:"memory_kb"`
	SubmittedAt time.Time `json:"submitted_at"`
	Error       string    `json:"error"`
}
