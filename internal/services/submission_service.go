package services

import (
	"errors"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
)

type SubmissionService struct {
	submissionRepo *repositories.ProblemSubmissionRepository
}

func NewSubmissionService(repo *repositories.ProblemSubmissionRepository) *SubmissionService {
	return &SubmissionService{
		submissionRepo: repo,
	}
}

// Submit creates a new submission for a problem
func (s *SubmissionService) Submit(userID, problemID int, code string) error {
	if code == "" {
		return errors.New("code cannot be empty")
	}

	submission := &entities.ProblemSubmission{
		UserID:    userID,
		ProblemID: problemID,
		Code:      code,
		Status:    "pending", // Initial status
	}

	_, err := s.submissionRepo.Create(submission)

	return err
}

// GetUserSubmission retrieves a specific submission for a user and problem
func (s *SubmissionService) GetUserSubmission(userID, problemID int) (*entities.ProblemSubmission, error) {
	return s.submissionRepo.GetByUserAndProblem(userID, problemID)
}

// GetAllUserSubmissions retrieves all submissions for a user
func (s *SubmissionService) GetAllUserSubmissions(userID int) ([]*entities.ProblemSubmission, error) {
	return s.submissionRepo.GettAllByUser(userID)
}

// UpdateSubmissionStatus updates the status and results of a submission
func (s *SubmissionService) UpdateSubmissionStatus(submission *entities.ProblemSubmission) error {
	_, err := s.submissionRepo.Create(submission) // Uses Create since it's upsert-style in the repo

	return err
}
