package repositories

import (
	"database/sql"
	"errors"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
)

type ProblemSubmissionRepository struct {
	db *sql.DB
}

func NewProblemSubmissionRepository(db *sql.DB) *ProblemSubmissionRepository {
	return &ProblemSubmissionRepository{db: db}
}

func (r *ProblemSubmissionRepository) Create(submission *entities.ProblemSubmission) error {
	query := `
		INSERT INTO problem_submissions (
			user_id, problem_id, code, status, runtime_ms, error
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query,
		submission.UserID,
		submission.ProblemID,
		submission.Code,
		submission.Status,
		submission.Runtime,
		submission.Error,
	)
	return err
}

func (r *ProblemSubmissionRepository) GetByUserAndProblem(userID, problemID uint64) (*entities.ProblemSubmission, error) {
	submission := &entities.ProblemSubmission{}
	query := `
		SELECT user_id, problem_id, code, status, runtime_ms, error
		FROM problem_submissions WHERE user_id = $1 AND problem_id = $2`

	err := r.db.QueryRow(query, userID, problemID).Scan(
		&submission.UserID,
		&submission.ProblemID,
		&submission.Code,
		&submission.Status,
		&submission.Runtime,
		&submission.Error,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("Submission not found")
	}
	if err != nil {
		return nil, err
	}

	return submission, nil
}

func (r *ProblemSubmissionRepository) GettAllByUser(userID uint64) ([]*entities.ProblemSubmission, error) {
	query := `
		SELECT user_id, problem_id, code, status, runtime_ms, error
		FROM problem_submissions 
		WHERE user_id = $1
		ORDER BY problem_id`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []*entities.ProblemSubmission
	for rows.Next() {
		submission := &entities.ProblemSubmission{}
		err := rows.Scan(
			&submission.UserID,
			&submission.ProblemID,
			&submission.Code,
			&submission.Status,
			&submission.Runtime,
			&submission.Error,
		)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}

	return submissions, rows.Err()
}
