package repositories

import (
	"database/sql"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
)

type ProblemTestCaseRepository interface {
	GetTestCasesByProblemID(problemID int) ([]*entities.ProblemTestCase, error)
}

type ProblemTestCaseRepositoryImpl struct {
	db *sql.DB
}

func NewProblemTestCaseRepository(db *sql.DB) ProblemTestCaseRepository {
	return &ProblemTestCaseRepositoryImpl{
		db: db,
	}
}

func (r *ProblemTestCaseRepositoryImpl) GetTestCasesByProblemID(problemID int) ([]*entities.ProblemTestCase, error) {
	query := `
        SELECT problem_id, order_index, input, output
        FROM problem_test_cases
        WHERE problem_id = $1
        ORDER BY order_index`

	rows, err := r.db.Query(query, problemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var testCases []*entities.ProblemTestCase
	for rows.Next() {
		tc := &entities.ProblemTestCase{}
		err := rows.Scan(&tc.ProblemID, &tc.OrderIndex, &tc.Input, &tc.Output)
		if err != nil {
			return nil, err
		}
		testCases = append(testCases, tc)
	}
	return testCases, rows.Err()
}
