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

type ProblemSolutionDisplay struct {
	Solution    *ProblemSolution `json:"solution"`
	SubmittedBy string           `json:"submittedBy"`
}

type ProblemSolutionVoteStatus struct {
	Upvoted   bool `json:"upvoted"`
	Downvoted bool `json:"downvoted"`
}
