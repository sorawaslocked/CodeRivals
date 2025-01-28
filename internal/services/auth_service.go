package services

import (
	"errors"
	"github.com/sorawaslocked/CodeRivals/internal/dtos"
	"github.com/sorawaslocked/CodeRivals/internal/repositories"
	v "github.com/sorawaslocked/CodeRivals/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(f *dtos.UserLoginForm) (uint64, error) {
	f.Validator.Check(v.NotBlank(f.Username), "username", "username should not be blank")
	f.Validator.Check(v.NotBlank(f.Password), "password", "password should not be blank")

	userFromDb, err := s.userRepo.GetByUsername(f.Username)

	if err != nil {
		f.Validator.AddFieldError("username", "username is invalid")

		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(userFromDb.HashedPassword, []byte(f.Password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			f.Validator.AddFieldError("password", "password is invalid")
		}

		return 0, err
	}

	return userFromDb.ID, nil
}
