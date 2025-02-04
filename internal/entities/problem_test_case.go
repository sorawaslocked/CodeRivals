package entities

type ProblemTestCase struct {
	ProblemID  int    `json:"problem_id"`
	OrderIndex int    `json:"order_index"`
	Input      string `json:"input"`
	Output     string `json:"output"`
}
