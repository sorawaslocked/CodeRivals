package repositories

import (
	"database/sql"
	"github.com/sorawaslocked/CodeRivals/internal/dtos"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
)

type ProblemRepository interface {
	Count() (uint64, error)
	Create(request *dtos.ProblemCreateRequest) error
	Get(id uint64) (*entities.Problem, error)
	GetAll() ([]*entities.Problem, error)
	Delete(id uint64) error
}

type PGProblemRepository struct {
	db *sql.DB
}

func NewPGProblemRepository(db *sql.DB) ProblemRepository {
	return &PGProblemRepository{db: db}
}

func (repo *PGProblemRepository) Count() (uint64, error) {
	var count uint64

	err := repo.db.QueryRow("SELECT COUNT(*) FROM problems").Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *PGProblemRepository) Create(request *dtos.ProblemCreateRequest) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	problemStmt := `INSERT INTO problems (title, description, difficulty, created_at, updated_at)
	VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
	var problemId uint64

	err = tx.QueryRow(problemStmt, request.Title, request.Description, request.Difficulty).Scan(&problemId)
	if err != nil {
		return err
	}

	for _, topic := range request.Topics {
		topicStmt := `INSERT INTO problem_topics (problem_id, topic_id)
		VALUES ($1, $2)`

		_, err = tx.Exec(topicStmt, problemId, topic.ID)

		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	return err
}

func (repo *PGProblemRepository) Get(id uint64) (*entities.Problem, error) {

	return nil, nil
}

func (repo *PGProblemRepository) GetAll() ([]*entities.Problem, error) {
	return nil, nil
}

func (repo *PGProblemRepository) Delete(id uint64) error {
	return nil
}
