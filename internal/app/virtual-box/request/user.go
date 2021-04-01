package request

import validation "github.com/go-ozzo/ozzo-validation"

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (r UserRequest) Validate() error {
	return validation.ValidateStruct(
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}
