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

func (s *ProblemService) GetProblemByURL(url string) (*entities.Problem, error) {
	return s.problemRepo.GetByURL(url)
}

func (s *ProblemService) GetPaginatedProblemsWithTopics(offset, itemsPerPage int, topicIDs []int) ([]*entities.Problem, int, error) {
	problems, err := s.GetAllProblems()
	if err != nil {
		return nil, 0, err
	}

	// If no topics selected, return all problems
	if len(topicIDs) == 0 {
		totalProblems := len(problems)
		if totalProblems < offset {
			return nil, totalProblems, nil
		}
		if totalProblems-itemsPerPage < offset {
			return problems[offset:], totalProblems, nil
		}
		return problems[offset : offset+itemsPerPage], totalProblems, nil
	}

	// Filter problems by topics
	var filteredProblems []*entities.Problem
	for _, problem := range problems {
		hasAllTopics := true
		for _, topicID := range topicIDs {
			found := false
			for _, problemTopic := range problem.Topics {
				if problemTopic.ID == topicID {
					found = true
					break
				}
			}
			if !found {
				hasAllTopics = false
				break
			}
		}
		if hasAllTopics {
			filteredProblems = append(filteredProblems, problem)
		}
	}

	totalProblems := len(filteredProblems)
	if totalProblems < offset {
		return nil, totalProblems, nil
	}
	if totalProblems-itemsPerPage < offset {
		return filteredProblems[offset:], totalProblems, nil
	}
	return filteredProblems[offset : offset+itemsPerPage], totalProblems, nil
}

func (s *ProblemService) CreateProblemSolution(solution *entities.ProblemSolution) error {
	return s.problemRepo.CreateProblemSolution(solution)
}

func (s *ProblemService) GetSolutionsForProblem(problemId int) ([]*entities.ProblemSolution, error) {
	return s.problemRepo.GetSolutionsForProblem(problemId)
}

func (s *ProblemService) GetSolutionById(id int) (*entities.ProblemSolution, error) {
	return s.problemRepo.GetSolutionById(id)
}

func (s *ProblemService) GetUpvoteBySolutionIdAndUserId(solutionId int, userId int) (bool, error) {
	return s.problemRepo.GetUpvoteBySolutionIdAndUserId(solutionId, userId)
}

func (s *ProblemService) UpvoteSolution(solutionId int, userId int) error {
	err := s.problemRepo.UpsertSolutionVote(solutionId, userId, true)

	if err != nil {
		return err
	}

	return nil
}

func (s *ProblemService) DownvoteSolution(solutionId int, userId int) error {
	err := s.problemRepo.UpsertSolutionVote(solutionId, userId, false)

	if err != nil {
		return err
	}

	return nil
}

func (s *ProblemService) UnvoteSolution(solutionId int, userId int) error {
	return s.problemRepo.RemoveSolutionVote(solutionId, userId)
}
