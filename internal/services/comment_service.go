package services

import (
	"github.com/sorawaslocked/CodeRivals/internal/entities"
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
)

type CommentService struct {
	repo repositories.CommentRepository
}

func NewCommentService(repo repositories.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(userID, problemID int, text string) error {
	comment := &entities.Comment{
		UserID:    userID,
		ProblemID: problemID,
		CommentID: nil,
		TextValue: text,
	}
	return s.repo.Create(comment)
}

func (s *CommentService) CreateReply(userID, problemID int, parentCommentID int, text string) error {
	comment := &entities.Comment{
		UserID:    userID,
		ProblemID: problemID,
		CommentID: &parentCommentID,
		TextValue: text,
	}
	return s.repo.Create(comment)
}

func (s *CommentService) GetProblemComments(problemID int) ([]entities.Comment, error) {
	return s.repo.GetByProblemID(problemID)
}

func (s *CommentService) DeleteComment(id int) error {
	return s.repo.Delete(id)
}

func (s *CommentService) UpdateComment(id int, text string) error {
	comment := &entities.Comment{
		ID:        id,
		TextValue: text,
	}
	return s.repo.Update(comment)
}
