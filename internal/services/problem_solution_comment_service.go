package services

import (
	"github.com/sorawaslocked/CodeRivals/internal/entities"
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
)

type SolutionCommentService struct {
	repo repositories.SolutionCommentRepository
}

func NewSolutionCommentService(repo repositories.SolutionCommentRepository) *SolutionCommentService {
	return &SolutionCommentService{repo: repo}
}

func (s *SolutionCommentService) CreateComment(userID, solutionID int, text string) error {
	comment := &entities.ProblemSolutionComment{
		UserID:     userID,
		SolutionID: solutionID,
		CommentID:  nil,
		TextValue:  text,
	}
	return s.repo.Create(comment)
}

func (s *SolutionCommentService) CreateReply(userID, solutionID int, text string, parentCommentID *int) error {
	comment := &entities.ProblemSolutionComment{
		UserID:     userID,
		SolutionID: solutionID,
		CommentID:  parentCommentID,
		TextValue:  text,
	}
	return s.repo.Create(comment)
}

func (s *SolutionCommentService) GetSolutionComments(solutionID int) ([]entities.ProblemSolutionComment, error) {
	return s.repo.GetBySolutionID(solutionID)
}

func (s *SolutionCommentService) DeleteComment(id int) error {
	return s.repo.Delete(id)
}

func (s *SolutionCommentService) UpdateComment(id int, text string) error {
	comment := &entities.ProblemSolutionComment{
		ID:        id,
		TextValue: text,
	}
	return s.repo.Update(comment)
}
