package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var passwordRule = []validation.Rule {
	validation.Required,
	validation.Length(8, 32),
	validation.Match(regexp.MustCompile(`[a-zA-Z\d]*[a-z][a-zA-Z\d]*[A-Z][a-zA-Z\d]*\d[a-zA-Z\d]*`)).Error("cannot contain whitespaces"),
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the RegisterRequest struct.
// It checks that all fields are filled in, that the email is a valid email address,
// and that the password is between 8 and 32 characters, and does not contain any whitespace.
func (a RegisterRequest) Validate() error {
	return validation.ValidateStruct(&a, 
		validation.Field(&a.Name, validation.Required, validation.Length(3, 64)),
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, passwordRule...),
	)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the LoginRequest struct.
// It checks that the email is a valid email address,
// and that the password is between 8 and 32 characters, and does not contain any whitespace.
func (a LoginRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, passwordRule...),
	)
}

type RefreshRequest struct {
	Token	string `json:"token"`
}

// Validate validates the RefreshRequest struct.
// It checks that the token is required and does not contain any whitespace.
func (a RefreshRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(
			&a.Token,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\S+$`)).Error("cannot contain whitespaces"),
		),
	)
}

type NoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Validate validates the NoteRequest struct.
// It checks that both the title and content are required.
func (a NoteRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Title, validation.Required),
		validation.Field(&a.Content, validation.Required),
	)
}
