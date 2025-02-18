package entities

type AdminStats struct {
	TotalUsers            int
	TotalProblems         int
	TotalSubmissions      int
	TotalComments         int
	SuccessfulSubmissions int
}

type ProblemSuccessRate struct {
	ProblemID   int
	Title       string
	Submissions int
	SuccessRate float64
}

type ActiveUser struct {
	UserID      int
	Username    string
	Submissions int
	Solutions   int
	Comments    int
}

type AdminUser struct {
	*User
	IsAdmin bool
}

type UserDetails struct {
	Submissions []*FullProblemSubmission  `json:"submissions"`
	Solutions   []*AdminUserSolution      `json:"solutions"`
	Comments    []*ProblemSolutionComment `json:"comments"`
}

type AdminUserSolution struct {
	Solution     *ProblemSolution `json:"solution"`
	ProblemTitle string           `json:"problem_title"`
}

type TopicWithCount struct {
	*Topic
	ProblemsCount int
}
