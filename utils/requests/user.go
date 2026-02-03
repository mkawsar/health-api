package requests

import validation "github.com/go-ozzo/ozzo-validation"

type UserRequest struct {
	Name string `json:"name"`
}

// Validate validates the UserRequest struct.
// It checks that the name is required and has a length between 3 and 64.
func (a UserRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required, validation.Length(3, 64)),
	)
}
