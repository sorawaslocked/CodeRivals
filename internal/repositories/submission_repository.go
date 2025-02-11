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

func (r *ProblemSubmissionRepository) Create(submission *entities.ProblemSubmission) (int, error) {
	query := `
		INSERT INTO problem_submissions (
			user_id, problem_id, code, status, runtime_ms, memory_kb, submitted_at, error
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7) RETURNING id`

	err := r.db.QueryRow(query,
		submission.UserID,
		submission.ProblemID,
		submission.Code,
		submission.Status,
		submission.Runtime,
		submission.Memory,
		submission.Error,
	).Scan(&submission.ID)

	return submission.ID, err
}

func (r *ProblemSubmissionRepository) GetByID(id int) (*entities.ProblemSubmission, error) {
	submission := &entities.ProblemSubmission{}
	submission.ID = id

	stmt := `SELECT user_id, problem_id, code, status, runtime_ms, memory_kb, submitted_at, error
	FROM problem_submissions WHERE id = $1`

	err := r.db.QueryRow(stmt, id).Scan(
		&submission.UserID,
		&submission.ProblemID,
		&submission.Code,
		&submission.Status,
		&submission.Runtime,
		&submission.Memory,
		&submission.SubmittedAt,
		&submission.Error)

	if err != nil {
		return nil, err
	}

	return submission, nil
}

func (r *ProblemSubmissionRepository) GetByUserAndProblem(userID, problemID int) (*entities.ProblemSubmission, error) {
	submission := &entities.ProblemSubmission{
		UserID:    userID,
		ProblemID: problemID,
	}
	query := `
		SELECT id, code, status, runtime_ms, memory_kb, submitted_at, error
		FROM problem_submissions WHERE user_id = $1 AND problem_id = $2`

	err := r.db.QueryRow(query, userID, problemID).Scan(
		&submission.ID,
		&submission.Code,
		&submission.Status,
		&submission.Runtime,
		&submission.Memory,
		&submission.SubmittedAt,
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

func (r *ProblemSubmissionRepository) GetAllByUser(userID int) ([]*entities.ProblemSubmission, error) {
	query := `
		SELECT id, problem_id, code, status, runtime_ms, memory_kb, submitted_at, error
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
		submission := &entities.ProblemSubmission{
			UserID: userID,
		}

		err := rows.Scan(
			&submission.ID,
			&submission.ProblemID,
			&submission.Code,
			&submission.Status,
			&submission.Runtime,
			&submission.Memory,
			&submission.SubmittedAt,
			&submission.Error,
		)

		if err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}

	return submissions, rows.Err()
}
