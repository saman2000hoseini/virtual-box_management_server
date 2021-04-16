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

type ChangeStateRequest struct {
	ID    string `json:"id"`
	State string `json:"state"`
}

type ExecRequest struct {
	ID       string  `json:"id"`
	Command  string  `json:"command"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type UploadRequest struct {
	ID              string `json:"id"`
	SrcPath         string `json:"src_path"`
	DestinationPath string `json:"destination_path"`
}

type TransferRequest struct {
	Src             string `json:"src"`
	Destination     string `json:"destination"`
	SrcPath         string `json:"src_path"`
	DestinationPath string `json:"destination_path"`
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

func (r ChangeStateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.State, validation.Required),
	)
}

func (r ExecRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.Command, validation.Required),
	)
}

func (r UploadRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.SrcPath, validation.Required),
		validation.Field(&r.DestinationPath, validation.Required),
	)
}

func (r TransferRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Src, validation.Required),
		validation.Field(&r.Destination, validation.Required),
		validation.Field(&r.SrcPath, validation.Required),
		validation.Field(&r.DestinationPath, validation.Required),
	)
}
