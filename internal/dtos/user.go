package dtos

import "github.com/sorawaslocked/CodeRivals/internal/validator"

type UserLoginForm struct {
	Username string
	Password string
	validator.Validator
}
