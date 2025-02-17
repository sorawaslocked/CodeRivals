package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
	"strings"
)

type SubmissionService struct {
	submissionRepo *repositories.ProblemSubmissionRepository
	problemService *ProblemService
	userRepo       repositories.UserRepository
	execService    *CodeExecutionService
	problemRepo    repositories.ProblemRepository
	testCaseRepo   repositories.ProblemTestCaseRepository
}

func NewSubmissionService(repo *repositories.ProblemSubmissionRepository, problemService *ProblemService, userRepo repositories.UserRepository, execService *CodeExecutionService, problemRepo repositories.ProblemRepository, testCaseRepo repositories.ProblemTestCaseRepository) *SubmissionService {
	return &SubmissionService{
		submissionRepo: repo,
		problemService: problemService,
		userRepo:       userRepo,
		execService:    execService,
		problemRepo:    problemRepo,
		testCaseRepo:   testCaseRepo,
	}
}

// Submit creates a new submission for a problem
func (s *SubmissionService) Submit(userID, problemID int, code string) (*entities.ProblemSubmission, error) {
	if code == "" {
		return nil, errors.New("code cannot be empty")
	}

	submission := &entities.ProblemSubmission{
		UserID:    userID,
		ProblemID: problemID,
		Code:      code,
	}

	problem, err := s.problemRepo.Get(problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	testCases, err := s.testCaseRepo.GetTestCasesByProblemID(problem.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test cases: %w", err)
	}

	err = s.ProcessSubmission(problem, testCases, submission)
	if err != nil {
		return nil, fmt.Errorf("failed to process submission: %w", err)
	}

	return submission, nil
}

func (s *SubmissionService) GetSubmission(submissionID int) (*entities.ProblemSubmission, error) {
	return s.submissionRepo.GetByID(submissionID)
}

// GetUserSubmission retrieves a specific submission for a user and problem
func (s *SubmissionService) GetUserSubmission(userID, problemID int) (*entities.ProblemSubmission, error) {
	return s.submissionRepo.GetByUserAndProblem(userID, problemID)
}

// GetAllUserSubmissions retrieves all submissions for a user
func (s *SubmissionService) GetAllUserSubmissions(userID int) ([]*entities.FullProblemSubmission, error) {
	submissions, err := s.submissionRepo.GetAllByUser(userID)

	if err != nil {
		return nil, err
	}

	var fullSubmissions []*entities.FullProblemSubmission

	for _, submission := range submissions {
		problem, err := s.problemService.GetProblem(submission.ProblemID)

		if err != nil {
			return nil, err
		}

		fullSubmissions = append(fullSubmissions, &entities.FullProblemSubmission{
			Submission: submission,
			Problem:    problem,
		})
	}

	return fullSubmissions, nil
}

// UpdateSubmissionStatus updates the status and results of a submission
func (s *SubmissionService) UpdateSubmissionStatus(submission *entities.ProblemSubmission) error {
	_, err := s.submissionRepo.Create(submission) // Uses Create since it's upsert-style in the repo

	return err
}

func (s *SubmissionService) ProcessSubmission(problem *entities.Problem, testCases []*entities.ProblemTestCase, submission *entities.ProblemSubmission) error {
	result := s.execService.ExecuteSolution(problem, testCases, submission)

	submission.Runtime = uint32(result.TimeMS)
	submission.Memory = uint32(result.MemoryKB)

	if result.Error != "" {
		if strings.Contains(result.Error, "runtime error") {
			submission.Status = "runtime_error"
			submission.Error = result.Error
		} else if strings.Contains(result.Error, "time limit") {
			submission.Status = "time_limit"
			submission.Error = result.Error
		} else {
			submission.Status = "compilation_error"
			submission.Error = result.Error
		}
	} else if result.Success {
		submission.Status = "accepted"

		previousSubmission, err := s.submissionRepo.GetByUserAndProblem(submission.UserID, submission.ProblemID)
		if err != nil {
			if err.Error() == "Submission not found" {
				points := getPointsForDifficulty(problem.Difficulty)
				if err := s.userRepo.AddPoints(submission.UserID, points); err != nil {
					return fmt.Errorf("failed to award points: %w", err)
				}
			} else {
				return fmt.Errorf("failed to check previous submissions: %w", err)
			}
		} else if previousSubmission.Status != "accepted" {
			points := getPointsForDifficulty(problem.Difficulty)
			if err := s.userRepo.AddPoints(submission.UserID, points); err != nil {
				return fmt.Errorf("failed to award points: %w", err)
			}
		}
	} else {
		submission.Status = "wrong_answer"
		for _, testResult := range result.TestResults {
			if !testResult.Passed {
				errorInfo := map[string]interface{}{
					"input":    testResult.Input,
					"expected": testResult.Expected,
					"output":   testResult.Actual,
				}
				errorJSON, err := json.Marshal(errorInfo)
				if err != nil {
					return fmt.Errorf("failed to marshal test results: %w", err)
				}
				submission.Error = string(errorJSON)
				break
			}
		}
	}
	_, err := s.submissionRepo.Create(submission)
	return err
}

func getPointsForDifficulty(difficulty string) int {
	switch difficulty {
	case "Easy":
		return 100
	case "Medium":
		return 300
	case "Hard":
		return 500
	default:
		return 0
	}
}
