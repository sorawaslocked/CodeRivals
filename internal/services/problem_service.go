package services

import (
	"errors"
	"github.com/sorawaslocked/CodeRivals/internal/dtos"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
)

type ProblemService struct {
	problemRepo repositories.ProblemRepository
}

func NewProblemService(repo repositories.ProblemRepository) *ProblemService {
	return &ProblemService{
		problemRepo: repo,
	}
}

func (s *ProblemService) GetProblem(id int) (*entities.Problem, error) {
	return s.problemRepo.Get(id)
}

func (s *ProblemService) GetAllProblems() ([]*entities.Problem, error) {
	return s.problemRepo.GetAll()
}

func (s *ProblemService) GetPaginatedProblems(offset, itemsPerPage int) ([]*entities.Problem, int, error) {
	problems, err := s.GetAllProblems()

	if err != nil {
		return nil, 0, err
	}

	totalProblems := len(problems)

	if totalProblems < offset {
		return nil, totalProblems, nil
	}

	if totalProblems-itemsPerPage < offset {
		return problems[offset:], totalProblems, nil
	}

	return problems[offset : offset+itemsPerPage], totalProblems, nil
}

func (s *ProblemService) ListProblems() ([]*entities.Problem, error) {
	return s.problemRepo.GetAll()
}

func (s *ProblemService) UpdateProblem(problem *entities.Problem) error {
	if problem.Title == "" {
		return errors.New("problem title cannot be empty")
	}

	if problem.Description == "" {
		return errors.New("problem description cannot be empty")
	}

	return s.problemRepo.Update(problem)
}

func (s *ProblemService) DeleteProblem(id int) error {
	return s.problemRepo.Delete(id)
}

func (s *ProblemService) CreateProblem(req *dtos.ProblemCreateRequest) error {
	if req.Title == "" {
		return errors.New("problem title cannot be empty")
	}

	if req.Description == "" {
		return errors.New("problem description cannot be empty")
	}

	if req.Difficulty == "" {
		return errors.New("problem difficulty cannot be empty")
	}

	problem := &entities.Problem{
		Title:       req.Title,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		Topics:      req.Topics,
	}

	return s.problemRepo.Create(problem)
}

func (s *ProblemService) GetProblemExamples(problemID int) ([]entities.ProblemExample, error) {
	return s.problemRepo.GetProblemExamples(problemID)
}
