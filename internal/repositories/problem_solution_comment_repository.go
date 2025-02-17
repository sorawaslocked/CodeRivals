package repositories

import (
	"database/sql"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
)

type SolutionCommentRepository interface {
	Create(comment *entities.ProblemSolutionComment) error
	GetBySolutionID(solutionID int) ([]entities.ProblemSolutionComment, error)
	Delete(id int) error
	Update(comment *entities.ProblemSolutionComment) error
}

type SolutionCommentRepositoryImpl struct {
	DB *sql.DB
}

func NewSolutionCommentRepository(db *sql.DB) *SolutionCommentRepositoryImpl {
	return &SolutionCommentRepositoryImpl{DB: db}
}

func (repo *SolutionCommentRepositoryImpl) Create(comment *entities.ProblemSolutionComment) error {
	var query string
	var err error

	if comment.CommentID == nil {
		query = `
            INSERT INTO problem_solution_comments (user_id, solution_id, text_value)
            VALUES ($1, $2, $3)
            RETURNING id, created_at`
		err = repo.DB.QueryRow(query,
			comment.UserID,
			comment.SolutionID,
			comment.TextValue).Scan(&comment.ID, &comment.CreatedAt)
	} else {
		query = `
            INSERT INTO problem_solution_comments (user_id, solution_id, comment_id, text_value)
            VALUES ($1, $2, $3, $4)
            RETURNING id, created_at`
		err = repo.DB.QueryRow(query,
			comment.UserID,
			comment.SolutionID,
			comment.CommentID,
			comment.TextValue).Scan(&comment.ID, &comment.CreatedAt)
	}

	return err
}

func (r *SolutionCommentRepositoryImpl) GetBySolutionID(solutionID int) ([]entities.ProblemSolutionComment, error) {
	query := `
        SELECT c.id, c.user_id, u.username, c.solution_id, c.comment_id, c.text_value, c.created_at
        FROM problem_solution_comments c
        JOIN users u on u.id = c.user_id
        WHERE solution_id = $1
        ORDER BY 
            CASE WHEN c.comment_id IS NULL THEN c.id END, 
            c.comment_id, 
            c.id`

	rows, err := r.DB.Query(query, solutionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []entities.ProblemSolutionComment
	for rows.Next() {
		var comment entities.ProblemSolutionComment
		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.Username,
			&comment.SolutionID,
			&comment.CommentID,
			&comment.TextValue,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *SolutionCommentRepositoryImpl) Delete(id int) error {
	query := `
    WITH RECURSIVE comment_tree AS (
        SELECT id FROM problem_solution_comments WHERE id = $1
        UNION
        SELECT c.id FROM problem_solution_comments c
        INNER JOIN comment_tree ct ON c.comment_id = ct.id
    )
    DELETE FROM problem_solution_comments WHERE id IN (SELECT id FROM comment_tree)`

	_, err := r.DB.Exec(query, id)
	return err
}

func (r *SolutionCommentRepositoryImpl) Update(comment *entities.ProblemSolutionComment) error {
	query := `
        UPDATE problem_solution_comments
        SET text_value = $1
        WHERE id = $2`

	_, err := r.DB.Exec(query, comment.TextValue, comment.ID)
	return err
}
