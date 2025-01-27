package entities

type ProblemSubmission struct {
	UserID    int    `json:"user_id"`
	ProblemID int    `json:"problem_id"`
	Code      string `json:"code"`
	Status    string `json:"status"`
	Runtime   int    `json:"runtime_ms"`
	Error     string `json:"error"`
}
