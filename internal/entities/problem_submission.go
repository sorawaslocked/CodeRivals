package entities

type ProblemSubmission struct {
	UserID    uint    `json:"user_id"`
	ProblemID uint    `json:"problem_id"`
	Code      string  `json:"code"`
	Status    string  `json:"status"`
	Runtime   int     `json:"runtime_ms"`
	Error     *string `json:"error"`
}
