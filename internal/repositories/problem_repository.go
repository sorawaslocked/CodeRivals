package repositories

import (
	"database/sql"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
)

type ProblemRepository interface {
	Count() (int, error)
	Create(request *entities.Problem) error
	Get(id int) (*entities.Problem, error)
	GetAll() ([]*entities.Problem, error)
	Update(problem *entities.Problem) error
	Delete(id int) error
	GetProblemExamples(problemID int) ([]entities.ProblemExample, error)
}

type PGProblemRepository struct {
	db              *sql.DB
	topicRepository TopicRepository
}

func NewPGProblemRepository(db *sql.DB, topicRepo TopicRepository) ProblemRepository {
	return &PGProblemRepository{
		db:              db,
		topicRepository: topicRepo,
	}
}

func (repo *PGProblemRepository) Count() (int, error) {
	var count int

	err := repo.db.QueryRow("SELECT COUNT(*) FROM problems").Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *PGProblemRepository) Create(problem *entities.Problem) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	problemStmt := `INSERT INTO problems (title, description, difficulty, url, created_at, updated_at)
	VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`
	var problemId int

	err = tx.QueryRow(problemStmt, problem.Title, problem.Description, problem.Difficulty, problem.Url).Scan(&problemId)
	if err != nil {
		return err
	}

	for _, topic := range problem.Topics {
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

func (repo *PGProblemRepository) Get(id int) (*entities.Problem, error) {
	prob := &entities.Problem{}

	probStmt := `SELECT title, description, difficulty, url, created_at, updated_at
	FROM problems WHERE id = $1`

	err := repo.db.QueryRow(probStmt, id).Scan(
		&prob.Title,
		&prob.Description,
		&prob.Difficulty,
		&prob.Url,
		&prob.CreatedAt,
		&prob.UpdatedAt)

	if err != nil {
		return nil, err
	}

	var topics []*entities.Topic
	topics, err = repo.topicRepository.GetAllForProblem(id)

	if err != nil {
		return nil, err
	}

	prob.Topics = topics

	return prob, nil
}

func (repo *PGProblemRepository) GetAll() ([]*entities.Problem, error) {
	probStmt := `SELECT id, title, description, difficulty, url, created_at, updated_at
	FROM problems`

	rows, err := repo.db.Query(probStmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var probs []*entities.Problem

	for rows.Next() {
		prob := &entities.Problem{}

		err = rows.Scan(
			&prob.ID,
			&prob.Title,
			&prob.Description,
			&prob.Difficulty,
			&prob.Url,
			&prob.CreatedAt,
			&prob.UpdatedAt)

		if err != nil {
			return nil, err
		}

		var topics []*entities.Topic
		topics, err = repo.topicRepository.GetAllForProblem(prob.ID)

		if err != nil {
			return nil, err
		}

		prob.Topics = topics
		probs = append(probs, prob)
	}

	return probs, nil
}

func (repo *PGProblemRepository) Update(prob *entities.Problem) error {
	tx, err := repo.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt := `UPDATE problems
	SET title = $1, description = $2, difficulty = $3, url = $4, updated_at = NOW()
	WHERE id = $5`

	_, err = tx.Exec(stmt, prob.Title, prob.Description, prob.Difficulty, prob.Url, prob.ID)

	if err != nil {
		return err
	}

	deleteTopicsStmt := `DELETE FROM problem_topics WHERE problem_id = $1`

	_, err = tx.Exec(deleteTopicsStmt, prob.ID)

	if err != nil {
		return err
	}

	for _, topic := range prob.Topics {
		insertTopicsStmt := `INSERT INTO problem_topics (problem_id, topic_id)
		VALUES ($1, $2)`

		_, err = tx.Exec(insertTopicsStmt, prob.ID, topic.ID)

		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (repo *PGProblemRepository) Delete(id int) error {
	tx, err := repo.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	topicsStmt := `DELETE FROM problem_topics WHERE problem_id = $1`

	_, err = tx.Exec(topicsStmt, id)

	if err != nil {
		return err
	}

	problemStmt := `DELETE FROM problems WHERE id = $1`

	_, err = tx.Exec(problemStmt, id)

	if err != nil {
		return err
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (r *PGProblemRepository) GetProblemExamples(problemID int) ([]entities.ProblemExample, error) {
	rows, err := r.db.Query(`
		SELECT given, expected, explanation
		FROM problem_examples
		WHERE problem_id = $1`, problemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var examples []entities.ProblemExample
	for rows.Next() {
		example := entities.ProblemExample{ProblemID: problemID}
		err := rows.Scan(&example.Given, &example.Expected, &example.Explanation)
		if err != nil {
			return nil, err
		}
		examples = append(examples, example)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return examples, nil
}
