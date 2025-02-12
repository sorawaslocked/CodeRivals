package services

import (
	"errors"
	"github.com/sorawaslocked/CodeRivals/internal/dtos"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
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

func (s *AuthService) GetUser(id int) (*entities.User, error) {
	return s.userRepo.Get(id)
}

func (s *AuthService) Login(f *dtos.UserLoginForm) (int, error) {
	f.Validator.Check(v.NotBlank(f.Username), "username", "Username should not be blank")

	userFromDb, err := s.userRepo.GetByUsername(f.Username)

	if err != nil {
		f.Validator.AddFieldError("username", "Username is invalid")

		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(userFromDb.HashedPassword, []byte(f.Password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			f.Validator.AddFieldError("password", "Password is invalid")
		}

		return 0, err
	}

	return userFromDb.ID, nil
}

func (s *AuthService) Register(f *dtos.UserRegisterForm) error {
	f.Validator.Check(v.NotBlank(f.Username), "username", "Username should not be blank")
	f.Validator.Check(v.NotBlank(f.Email), "email", "Email should not be blank")
	f.Validator.Check(v.MinChars(f.Password, 8), "password", "Password should not be less than 8 characters")
	f.Validator.Check(v.MaxChars(f.Password, 20), "password", "Password should not be more than 20 characters")
	f.Validator.Check(v.Equal(f.Password, f.ConfirmPassword), "password", "Passwords should be equal")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(f.Password), 10)

	if err != nil {
		return err
	}

	err = s.userRepo.Create(&entities.User{
		Username:       f.Username,
		Email:          f.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		if errors.Is(err, repositories.ErrDuplicateUsername) {
			f.Validator.AddFieldError("username", "Username already in use")

			return nil
		}
		if errors.Is(err, repositories.ErrDuplicateEmail) {
			f.Validator.AddFieldError("email", "Email already in use")

			return nil
		}

		return err
	}

	return nil
}
