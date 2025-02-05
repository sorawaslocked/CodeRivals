package repositories

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
)

type ProblemRepository interface {
	Count() (int, error)
	Create(request *entities.Problem) error
	Get(id int) (*entities.Problem, error)
	GetAll() ([]*entities.Problem, error)
	GetTestCases(problemId int) ([]*entities.ProblemTestCase, error)
	Update(problem *entities.Problem) error
	Delete(id int) error
	GetProblemExamples(problemID int) ([]entities.ProblemExample, error)
	GetByURL(url string) (*entities.Problem, error)
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

func (repo *PGProblemRepository) Create(prob *entities.Problem) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	problemStmt := `INSERT INTO problems
    (title, description, difficulty, url, input_types, output_type, method_name, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) RETURNING id`
	var problemId int

	err = tx.QueryRow(
		problemStmt,
		prob.Title,
		prob.Description,
		prob.Difficulty,
		prob.Url,
		prob.InputTypes,
		prob.OutputType,
		prob.MethodName,
	).Scan(&problemId)

	if err != nil {
		return err
	}

	for _, topic := range prob.Topics {
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
	prob.ID = id

	probStmt := `SELECT title, description, difficulty, url, input_types, output_type, method_name, created_at, updated_at
	FROM problems WHERE id = $1`

	err := repo.db.QueryRow(probStmt, id).Scan(
		&prob.Title,
		&prob.Description,
		&prob.Difficulty,
		&prob.Url,
		pq.Array(&prob.InputTypes),
		&prob.OutputType,
		&prob.MethodName,
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
	probStmt := `SELECT id, title, description, difficulty, url, input_types, output_type, method_name, created_at, updated_at
	FROM problems ORDER BY id`

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
			pq.Array(&prob.InputTypes),
			&prob.OutputType,
			&prob.MethodName,
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
	SET title = $1, description = $2, difficulty = $3, url = $4, input_types = $5, output_type = $6, method_name = $7, updated_at = NOW()
	WHERE id = $8`

	_, err = tx.Exec(
		stmt,
		prob.Title,
		prob.Description,
		prob.Difficulty,
		prob.Url,
		prob.InputTypes,
		prob.OutputType,
		prob.MethodName,
		prob.ID,
	)

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
		SELECT input, output, explanation
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

func (repo *PGProblemRepository) GetTestCases(problemId int) ([]*entities.ProblemTestCase, error) {
	stmt := `SELECT order_index, input, output
	FROM problem_test_cases
	WHERE problem_id = $1
	ORDER BY order_index`

	rows, err := repo.db.Query(stmt, problemId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var testCases []*entities.ProblemTestCase

	for rows.Next() {
		tc := &entities.ProblemTestCase{}
		tc.ProblemID = problemId

		err = rows.Scan(
			&tc.OrderIndex,
			&tc.Input,
			&tc.Output)

		if err != nil {
			return nil, err
		}

		testCases = append(testCases, tc)
	}

	return testCases, nil
}

func (repo *PGProblemRepository) GetByURL(url string) (*entities.Problem, error) {
	prob := &entities.Problem{}

	probStmt := `SELECT id, title, description, difficulty, url, input_types, output_type, method_name, created_at, updated_at 
    FROM problems WHERE url = $1`

	err := repo.db.QueryRow(probStmt, url).Scan(
		&prob.ID,
		&prob.Title,
		&prob.Description,
		&prob.Difficulty,
		&prob.Url,
		pq.Array(&prob.InputTypes),
		&prob.OutputType,
		&prob.MethodName,
		&prob.CreatedAt,
		&prob.UpdatedAt)

	if err != nil {
		return nil, err
	}

	topics, err := repo.topicRepository.GetAllForProblem(prob.ID)
	if err != nil {
		return nil, err
	}

	prob.Topics = topics
	return prob, nil
}
