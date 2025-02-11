package entities

type ProblemSolution struct {
	ID          int    `json:"id"`
	ProblemId   int    `json:"problemId"`
	UserId      int    `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Code        string `json:"code"`
	Votes       int    `json:"votes"`
}
