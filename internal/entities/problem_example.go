package entities

type ProblemExample struct {
	ProblemID   int                    `json:"problem_id"`
	Given       map[string]interface{} `json:"given"`
	Expected    map[string]interface{} `json:"expected"`
	Explanation string                 `json:"explanation"`
}
