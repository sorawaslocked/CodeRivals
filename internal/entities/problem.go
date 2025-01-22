package entities

import "time"

type Problem struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Difficulty  int       `json:"difficulty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Topics      []Topic   `json:"topics"`
}
