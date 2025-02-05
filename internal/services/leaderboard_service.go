package services

import (
	"github.com/sorawaslocked/CodeRivals/internal/entities"
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
)

type LeaderboardService struct {
	userRepo repositories.UserRepository
}

func NewLeaderboardService(userRepo repositories.UserRepository) *LeaderboardService {
	return &LeaderboardService{
		userRepo: userRepo,
	}
}

func (s *LeaderboardService) GetLeaderboard() ([]*entities.User, error) {
	return s.userRepo.GetUserRankings()
}
