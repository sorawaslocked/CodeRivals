package repositories

import (
	"database/sql"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
)

type CommentRepository interface {
	Create(comment *entities.Comment) error
	GetByProblemID(problemID int) ([]entities.Comment, error)
	Delete(id int) error
	Update(comment *entities.Comment) error
}

type CommentRepositoryImpl struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepositoryImpl {
	return &CommentRepositoryImpl{DB: db}
}

func (repo *CommentRepositoryImpl) Create(comment *entities.Comment) error {
	var query string
	var err error

	if comment.CommentID == nil {
		// For top-level comments
		query = `
            INSERT INTO comments (user_id, problem_id, text_value)
            VALUES ($1, $2, $3)
            RETURNING id`
		err = repo.DB.QueryRow(query,
			comment.UserID,
			comment.ProblemID,
			comment.TextValue).Scan(&comment.ID)
	} else {
		// For replies
		query = `
            INSERT INTO comments (user_id, problem_id, comment_id, text_value)
            VALUES ($1, $2, $3, $4)
            RETURNING id`
		err = repo.DB.QueryRow(query,
			comment.UserID,
			comment.ProblemID,
			comment.CommentID,
			comment.TextValue).Scan(&comment.ID)
	}

	return err
}

func (r *CommentRepositoryImpl) GetByProblemID(problemID int) ([]entities.Comment, error) {
	query := `
        SELECT c.id, c.user_id, u.username, c.problem_id, c.comment_id, c.text_value
        FROM comments c
        JOIN users u on u.id = c.user_id
        WHERE problem_id = $1
        ORDER BY 
            CASE WHEN c.comment_id IS NULL THEN c.id END, 
            c.comment_id, 
            c.id`

	rows, err := r.DB.Query(query, problemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []entities.Comment
	for rows.Next() {
		var comment entities.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.Username,
			&comment.ProblemID,
			&comment.CommentID,
			&comment.TextValue,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *CommentRepositoryImpl) Delete(id int) error {
	query := `
    WITH RECURSIVE comment_tree AS (
        SELECT id FROM comments WHERE id = $1
        UNION
        SELECT c.id FROM comments c
        INNER JOIN comment_tree ct ON c.comment_id = ct.id
    )
    DELETE FROM comments WHERE id IN (SELECT id FROM comment_tree)`

	_, err := r.DB.Exec(query, id)
	return err
}

func (r *CommentRepositoryImpl) Update(comment *entities.Comment) error {
	query := `
        UPDATE comments
        SET text_value = $1
        WHERE id = $2`

	_, err := r.DB.Exec(query, comment.TextValue, comment.ID)
	return err
}
