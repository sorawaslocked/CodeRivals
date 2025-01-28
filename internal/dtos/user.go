package dtos

import "github.com/sorawaslocked/CodeRivals/internal/validator"

type UserLoginForm struct {
	Username string
	Password string
	validator.Validator
}

type UserRegisterForm struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
	validator.Validator
}
