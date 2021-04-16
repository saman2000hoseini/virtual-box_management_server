package request

import validation "github.com/go-ozzo/ozzo-validation"

type CloneRequest struct {
	ID    string `json:"id"`
	NewID string `json:"new_id"`
}

type ModifyRequest struct {
	ID     string `json:"id"`
	CPUs   uint   `json:"cpus"`
	Memory uint   `json:"memory"`
}

type ExecRequest struct {
	ID      string `json:"id"`
	Command string `json:"command"`
}

func (r CloneRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.NewID, validation.Required),
	)
}

func (r ModifyRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
	)
}

func (r ExecRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.Command, validation.Required),
	)
}
