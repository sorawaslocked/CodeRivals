package entities

type Comment struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ProblemID int    `json:"problem_id"`
	CommentID *int   `json:"comment_id"`
	TextValue string `json:"text_value"`
	Username  string `json:"username"`
}
