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

func (s *ProblemService) GetProblem(id uint64) (*entities.Problem, error) {
	return s.problemRepo.Get(id)
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

func (s *ProblemService) DeleteProblem(id uint64) error {
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
