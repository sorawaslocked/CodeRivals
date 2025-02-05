package entities

import "encoding/json"

type ProblemExample struct {
	ProblemID   int             `json:"problem_id"`
	Given       json.RawMessage `json:"given"`
	Expected    json.RawMessage `json:"expected"`
	Explanation string          `json:"explanation"`
}
