package services

import (
	"github.com/sorawaslocked/CodeRivals/internal/entities"
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

type AdminService struct {
	roleRepo       repositories.RoleRepository
	userRepo       repositories.UserRepository
	problemRepo    repositories.ProblemRepository
	topicRepo      repositories.TopicRepository
	submissionRepo *repositories.ProblemSubmissionRepository
	commentRepo    *repositories.CommentRepositoryImpl
}

func NewAdminService(
	roleRepo repositories.RoleRepository,
	userRepo repositories.UserRepository,
	problemRepo repositories.ProblemRepository,
	topicRepo repositories.TopicRepository,
	submissionRepo *repositories.ProblemSubmissionRepository,
	commentRepo *repositories.CommentRepositoryImpl,
) *AdminService {
	return &AdminService{
		roleRepo:       roleRepo,
		userRepo:       userRepo,
		problemRepo:    problemRepo,
		topicRepo:      topicRepo,
		submissionRepo: submissionRepo,
		commentRepo:    commentRepo,
	}
}

func (s *AdminService) IsUserAdmin(userId int) (bool, error) {
	return s.roleRepo.IsUserAdmin(userId)
}

func (s *AdminService) GetStats() (map[string]int, error) {
	stats := make(map[string]int)

	userCount, err := s.userRepo.Count()
	if err != nil {
		return nil, err
	}
	stats["totalUsers"] = userCount

	problemCount, err := s.problemRepo.Count()
	if err != nil {
		return nil, err
	}
	stats["totalProblems"] = problemCount

	topicCount, err := s.topicRepo.Count()
	if err != nil {
		return nil, err
	}
	stats["totalTopics"] = topicCount

	return stats, nil
}

func (s *AdminService) GetAllUsersWithDetails() ([]*entities.AdminUser, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var adminUsers []*entities.AdminUser
	for _, user := range users {
		isAdmin, err := s.roleRepo.IsUserAdmin(user.ID)
		if err != nil {
			return nil, err
		}

		adminUser := &entities.AdminUser{
			User:    user,
			IsAdmin: isAdmin,
		}
		adminUsers = append(adminUsers, adminUser)
	}

	return adminUsers, nil
}

func (s *AdminService) GetUserDetails(userId int) (*entities.UserDetails, error) {
	// Get submissions
	submissions, err := s.submissionRepo.GetAllByUser(userId)
	if err != nil {
		return nil, err
	}

	var fullSubmissions []*entities.FullProblemSubmission
	for _, submission := range submissions {
		problem, err := s.problemRepo.Get(submission.ProblemID)
		if err != nil {
			return nil, err
		}
		fullSubmissions = append(fullSubmissions, &entities.FullProblemSubmission{
			Submission: submission,
			Problem:    problem,
		})
	}

	// Get solutions
	solutions, err := s.problemRepo.GetSolutionsByUser(userId)
	if err != nil {
		return nil, err
	}

	var adminSolutions []*entities.AdminUserSolution
	for _, solution := range solutions {
		problem, err := s.problemRepo.Get(solution.ProblemId)
		if err != nil {
			return nil, err
		}

		adminSolution := &entities.AdminUserSolution{
			Solution:     solution,
			ProblemTitle: problem.Title,
		}
		adminSolutions = append(adminSolutions, adminSolution)
	}

	// Get comments
	comments, err := s.commentRepo.GetByUser(userId)
	if err != nil {
		return nil, err
	}

	return &entities.UserDetails{
		Submissions: fullSubmissions,
		Solutions:   adminSolutions,
		Comments:    comments,
	}, nil
}

func (s *AdminService) UpdateUserRole(userId int, isAdmin bool) error {
	// Get admin role ID (we know it's 1 from initial DB setup)
	adminRoleId := 1

	if isAdmin {
		return s.roleRepo.AssignRole(userId, adminRoleId)
	}
	return s.roleRepo.RemoveRole(userId, adminRoleId)
}

func (s *AdminService) ResetUserPassword(userId int) (string, error) {
	// Generate a random password
	newPassword := generateRandomPassword()

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Update the user's password
	if err := s.userRepo.UpdatePassword(userId, hashedPassword); err != nil {
		return "", err
	}

	return newPassword, nil
}

func generateRandomPassword() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 12
	password := make([]byte, length)
	for i := range password {
		password[i] = chars[rand.Intn(len(chars))]
	}
	return string(password)
}
