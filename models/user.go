package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
)

// User user
// swagger:model User
type User struct {

	// password
	Password string `json:"password,omitempty"`

	// roles
	Roles []string `json:"roles,omitempty"`

	// username
	Username string `json:"username,omitempty"`
}

// Validate validates this user
func (m *User) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
