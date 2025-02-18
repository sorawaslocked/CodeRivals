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
	CreateProblemSolution(solution *entities.ProblemSolution) error
	GetSolutionsForProblem(problemId int) ([]*entities.ProblemSolution, error)
	GetSolutionById(id int) (*entities.ProblemSolution, error)
	GetSolutionsByUser(userId int) ([]*entities.ProblemSolution, error)
	GetVotesBySolutionId(solutionId int) (int, error)
	GetUpvoteBySolutionIdAndUserId(solutionId int, userId int) (bool, error)
	UpsertSolutionVote(solutionId int, userId int, upvote bool) error
	RemoveSolutionVote(solutionId int, userId int) error
	DeleteTestCases(problemId int) error
	CreateTestCases(testCases []*entities.ProblemTestCase) error
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
		pq.StringArray(prob.InputTypes),
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
		pq.StringArray(prob.InputTypes),
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

func (repo *PGProblemRepository) CreateProblemSolution(solution *entities.ProblemSolution) error {
	stmt := `INSERT INTO problem_solutions (problem_id, user_id, title, description, code)
	VALUES ($1, $2, $3, $4, $5)`

	_, err := repo.db.Exec(stmt,
		solution.ProblemId,
		solution.UserId,
		solution.Title,
		solution.Description,
		solution.Code)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PGProblemRepository) GetSolutionsForProblem(problemId int) ([]*entities.ProblemSolution, error) {
	stmt := `SELECT id, user_id, title, description, code FROM problem_solutions
	WHERE problem_id = $1`

	rows, err := repo.db.Query(stmt, problemId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var solutions []*entities.ProblemSolution

	for rows.Next() {
		sol := &entities.ProblemSolution{}
		sol.ProblemId = problemId

		err = rows.Scan(
			&sol.ID,
			&sol.UserId,
			&sol.Title,
			&sol.Description,
			&sol.Code)

		if err != nil {
			return nil, err
		}

		var votes int
		votes, err = repo.GetVotesBySolutionId(sol.ID)

		if err != nil {
			return nil, err
		}

		sol.Votes = votes

		solutions = append(solutions, sol)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return solutions, nil
}

func (repo *PGProblemRepository) GetSolutionById(id int) (*entities.ProblemSolution, error) {
	stmt := `SELECT problem_id, user_id, title, description, code FROM problem_solutions
	WHERE id = $1`

	sol := &entities.ProblemSolution{}
	sol.ID = id

	err := repo.db.QueryRow(stmt, id).Scan(
		&sol.ProblemId,
		&sol.UserId,
		&sol.Title,
		&sol.Description,
		&sol.Code)

	if err != nil {
		return nil, err
	}

	var votes int
	votes, err = repo.GetVotesBySolutionId(sol.ID)

	if err != nil {
		return nil, err
	}

	sol.Votes = votes

	return sol, nil
}

func (repo *PGProblemRepository) GetVotesBySolutionId(solutionId int) (int, error) {
	stmt := `SELECT COUNT(*)
 	FROM problem_solution_votes
 	WHERE solution_id = $1 AND upvote`

	var posCount int
	err := repo.db.QueryRow(stmt, solutionId).Scan(&posCount)

	if err != nil {
		return 0, err
	}

	stmt = `SELECT COUNT(*)
	FROM problem_solution_votes
	WHERE solution_id = $1 AND NOT upvote`

	var negCount int
	err = repo.db.QueryRow(stmt, solutionId).Scan(&negCount)

	if err != nil {
		return 0, err
	}

	return posCount - negCount, nil
}

func (repo *PGProblemRepository) GetUpvoteBySolutionIdAndUserId(solutionId int, userId int) (bool, error) {
	stmt := `SELECT upvote
	FROM problem_solution_votes
	WHERE solution_id = $1 AND user_id = $2`

	var upvote bool
	err := repo.db.QueryRow(stmt, solutionId, userId).Scan(&upvote)

	if err != nil {
		return false, err
	}

	return upvote, nil
}

func (repo *PGProblemRepository) UpsertSolutionVote(solutionId int, userId int, upvote bool) error {
	_, err := repo.GetUpvoteBySolutionIdAndUserId(solutionId, userId)

	if err == nil {
		stmt := `UPDATE problem_solution_votes
		SET upvote = $1
		WHERE solution_id = $2 AND user_id = $3`

		_, err = repo.db.Exec(stmt, upvote, solutionId, userId)

		if err != nil {
			return err
		}

		return nil
	}

	stmt := `INSERT INTO problem_solution_votes (solution_id, user_id, upvote)
	VALUES ($1, $2, $3)`

	_, err = repo.db.Exec(stmt, solutionId, userId, upvote)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PGProblemRepository) RemoveSolutionVote(solutionId int, userId int) error {
	stmt := `DELETE FROM problem_solution_votes
	WHERE solution_id = $1 AND user_id = $2`

	_, err := repo.db.Exec(stmt, solutionId, userId)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PGProblemRepository) GetSolutionsByUser(userId int) ([]*entities.ProblemSolution, error) {
	stmt := `SELECT id, problem_id, title, description, code 
            FROM problem_solutions 
            WHERE user_id = $1`

	rows, err := repo.db.Query(stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var solutions []*entities.ProblemSolution
	for rows.Next() {
		solution := &entities.ProblemSolution{UserId: userId}
		err := rows.Scan(
			&solution.ID,
			&solution.ProblemId,
			&solution.Title,
			&solution.Description,
			&solution.Code,
		)
		if err != nil {
			return nil, err
		}

		votes, err := repo.GetVotesBySolutionId(solution.ID)
		if err != nil {
			return nil, err
		}
		solution.Votes = votes

		solutions = append(solutions, solution)
	}

	return solutions, rows.Err()
}

func (repo *PGProblemRepository) DeleteTestCases(problemId int) error {
	stmt := `DELETE FROM problem_test_cases WHERE problem_id = $1`
	_, err := repo.db.Exec(stmt, problemId)
	return err
}

func (repo *PGProblemRepository) CreateTestCases(testCases []*entities.ProblemTestCase) error {
	stmt := `INSERT INTO problem_test_cases (problem_id, order_index, input, output)
             VALUES ($1, $2, $3, $4)`

	for _, tc := range testCases {
		_, err := repo.db.Exec(stmt, tc.ProblemID, tc.OrderIndex, tc.Input, tc.Output)
		if err != nil {
			return err
		}
	}
	return nil
}
